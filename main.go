package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer implements plugins for finding unused interface methods.
var Analyzer = &analysis.Analyzer{
	Name:     "unused_interface_methods",
	Doc:      "Checks for unused interface methods",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// methodInfo represents information about a method in an interface.
type methodInfo struct {
	ifaceName string           // interface name
	iface     *types.Interface // interface object
	method    *types.Func      // method object
	used      bool             // used flag
}

// collectInterfaceMethods collects all explicit interface methods in the package using parallel processing.
func collectInterfaceMethods(pass *analysis.Pass) map[*types.Func]methodInfo {
	var ifaceMethods sync.Map
	var wg sync.WaitGroup

	for _, file := range pass.Files {
		wg.Add(1)
		go func(file *ast.File) {
			defer wg.Done()

			for _, decl := range file.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.TYPE {
					continue
				}
				for _, spec := range gd.Specs {
					tspec := spec.(*ast.TypeSpec)
					if _, ok := tspec.Type.(*ast.InterfaceType); !ok {
						continue
					}
					obj := pass.TypesInfo.Defs[tspec.Name]
					if obj == nil {
						continue
					}
					named, ok := obj.Type().(*types.Named)
					if !ok {
						continue
					}
					ifaceType, ok := named.Underlying().(*types.Interface)
					if !ok {
						continue
					}
					for i := 0; i < ifaceType.NumExplicitMethods(); i++ {
						m := ifaceType.ExplicitMethod(i)
						if m == nil {
							continue
						}
						ifaceMethods.Store(m, methodInfo{
							ifaceName: tspec.Name.Name,
							iface:     ifaceType,
							method:    m,
							used:      false,
						})
					}
				}
			}
		}(file)
	}

	wg.Wait()

	// convert sync.Map to regular map
	result := make(map[*types.Func]methodInfo)
	ifaceMethods.Range(func(key, value interface{}) bool {
		result[key.(*types.Func)] = value.(methodInfo)
		return true
	})
	return result
}

// methodAnalyzer handles analysis of method usage in AST
type methodAnalyzer struct {
	pass         *analysis.Pass
	ifaceMethods map[*types.Func]methodInfo
	usedMethods  map[*types.Func]bool
	methodIndex  map[string][]*types.Func // index methods by name for faster lookup
}

// newMethodAnalyzer creates a new method analyzer
func newMethodAnalyzer(pass *analysis.Pass, ifaceMethods map[*types.Func]methodInfo) *methodAnalyzer {
	// build method index by name
	methodIndex := make(map[string][]*types.Func)
	for method := range ifaceMethods {
		name := method.Name()
		methodIndex[name] = append(methodIndex[name], method)
	}

	return &methodAnalyzer{
		pass:         pass,
		ifaceMethods: ifaceMethods,
		usedMethods:  make(map[*types.Func]bool),
		methodIndex:  methodIndex,
	}
}

// analyzeUsedMethods traverses AST and marks used methods
func analyzeUsedMethods(pass *analysis.Pass, ifaceMethods map[*types.Func]methodInfo) map[*types.Func]bool {
	analyzer := newMethodAnalyzer(pass, ifaceMethods)
	return analyzer.analyze()
}

// analyze performs the main analysis logic
func (ma *methodAnalyzer) analyze() map[*types.Func]bool {
	ins := ma.pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.SelectorExpr)(nil),
		(*ast.CallExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			ma.analyzeSelectorExpr(node)
		case *ast.CallExpr:
			ma.analyzeCallExpr(node)
		}
	})

	return ma.usedMethods
}

// analyzeSelectorExpr handles method calls through selectors
func (ma *methodAnalyzer) analyzeSelectorExpr(node *ast.SelectorExpr) {
	sel := ma.pass.TypesInfo.Selections[node]
	if sel == nil || (sel.Kind() != types.MethodVal && sel.Kind() != types.MethodExpr) {
		return
	}

	calledMethod := sel.Obj().(*types.Func)
	recv := sel.Recv()

	ma.markMatchingMethods(calledMethod, recv)
}

// markMatchingMethods marks interface methods that match the called method
func (ma *methodAnalyzer) markMatchingMethods(calledMethod *types.Func, recv types.Type) {
	// only check methods with matching names using index
	candidateMethods := ma.methodIndex[calledMethod.Name()]

	for _, ifaceMethod := range candidateMethods {
		if ma.usedMethods[ifaceMethod] {
			continue
		}

		info := ma.ifaceMethods[ifaceMethod]
		if ma.isMethodMatch(calledMethod, ifaceMethod, recv, info) {
			ma.usedMethods[ifaceMethod] = true
		}
	}
}

// isMethodMatch checks if called method matches interface method
func (ma *methodAnalyzer) isMethodMatch(calledMethod, ifaceMethod *types.Func, recv types.Type, info methodInfo) bool {
	// direct match
	if calledMethod == ifaceMethod {
		return true
	}

	// name mismatch
	if calledMethod.Name() != ifaceMethod.Name() {
		return false
	}

	// check if receiver implements interface
	return ma.checkImplements(recv, info)
}

// checkImplements checks if receiver type implements the interface
func (ma *methodAnalyzer) checkImplements(recv types.Type, info methodInfo) bool {
	// direct implementation
	if types.Implements(recv, info.iface) {
		return true
	}

	// interface-to-interface check
	if recvIface, ok := recv.Underlying().(*types.Interface); ok {
		if types.Implements(info.iface, recvIface) {
			return true
		}
	}

	// generic type check
	return ma.checkGenericImplements(recv, info)
}

// checkGenericImplements handles generic type implementations
func (ma *methodAnalyzer) checkGenericImplements(recv types.Type, info methodInfo) bool {
	namedRecv, ok := recv.(*types.Named)
	if !ok {
		return false
	}

	origin := namedRecv.Origin()
	if origin == nil {
		return false
	}

	ifaceOrig, ok := origin.Underlying().(*types.Interface)
	return ok && types.Identical(ifaceOrig, info.iface)
}

// analyzeCallExpr handles function calls (specifically fmt.* functions)
func (ma *methodAnalyzer) analyzeCallExpr(node *ast.CallExpr) {
	ident := ma.extractFunctionIdent(node)
	if ident == nil {
		return
	}

	if !ma.isFmtFunction(ident) {
		return
	}

	ma.analyzeFmtCall(node)
}

// extractFunctionIdent extracts function identifier from call expression
func (ma *methodAnalyzer) extractFunctionIdent(node *ast.CallExpr) *ast.Ident {
	switch fun := node.Fun.(type) {
	case *ast.Ident:
		return fun
	case *ast.SelectorExpr:
		return fun.Sel
	default:
		return nil
	}
}

// isFmtFunction checks if the function belongs to fmt package
func (ma *methodAnalyzer) isFmtFunction(ident *ast.Ident) bool {
	fn, ok := ma.pass.TypesInfo.Uses[ident].(*types.Func)
	return ok && fn.Pkg() != nil && fn.Pkg().Path() == "fmt"
}

// analyzeFmtCall analyzes fmt function calls for Stringer interface usage
func (ma *methodAnalyzer) analyzeFmtCall(node *ast.CallExpr) {
	for _, arg := range node.Args {
		argType := ma.pass.TypesInfo.TypeOf(arg)
		if argType == nil {
			continue
		}

		ma.checkStringerUsage(argType)
	}
}

// checkStringerUsage checks if argument implements Stringer interface
func (ma *methodAnalyzer) checkStringerUsage(argType types.Type) {
	// only check methods named "String" using index
	stringMethods := ma.methodIndex["String"]

	for _, ifaceMethod := range stringMethods {
		if ma.usedMethods[ifaceMethod] {
			continue
		}

		info := ma.ifaceMethods[ifaceMethod]
		if ma.isStringerMethod(ifaceMethod) && types.Implements(argType, info.iface) {
			ma.usedMethods[ifaceMethod] = true
		}
	}
}

// isStringerMethod checks if method is String() string
func (ma *methodAnalyzer) isStringerMethod(method *types.Func) bool {
	if method.Name() != "String" {
		return false
	}

	sig, ok := method.Type().(*types.Signature)
	if !ok || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}

	basic, ok := sig.Results().At(0).Type().(*types.Basic)
	return ok && basic.Kind() == types.String
}

// reportUnusedMethods sorts and reports methods that were not used.
func reportUnusedMethods(pass *analysis.Pass, ifaceMethods map[*types.Func]methodInfo, used map[*types.Func]bool) {
	// mark used methods
	for m := range used {
		if info, ok := ifaceMethods[m]; ok {
			info.used = true
			ifaceMethods[m] = info
		}
	}

	var unused []methodInfo
	for _, info := range ifaceMethods {
		if !info.used {
			unused = append(unused, info)
		}
	}

	sort.Slice(unused, func(i, j int) bool {
		posI := pass.Fset.Position(unused[i].method.Pos())
		posJ := pass.Fset.Position(unused[j].method.Pos())
		if posI.Filename != posJ.Filename {
			return posI.Filename < posJ.Filename
		}
		return posI.Line < posJ.Line
	})

	for _, info := range unused {
		pass.Reportf(info.method.Pos(), "method %q of interface %q is declared but not used", info.method.Name(), info.ifaceName)
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	ifaceMethods := collectInterfaceMethods(pass)
	used := analyzeUsedMethods(pass, ifaceMethods)
	reportUnusedMethods(pass, ifaceMethods, used)
	return nil, nil
}

func main() {
	singlechecker.Main(Analyzer)
}
