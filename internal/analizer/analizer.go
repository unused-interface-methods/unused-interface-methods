package analizer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

// a implements plugin for finding unused interface methods.
var a = &analysis.Analyzer{
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

// collectInterfaceMethods collects all explicit interface methods in the package.
func collectInterfaceMethods(pass *analysis.Pass) map[*types.Func]methodInfo {
	ifaceMethods := map[*types.Func]methodInfo{}

	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		relPath, err := filepath.Rel(basePath, filename)
		if err != nil {
			relPath = filename
		}

		if cfg.ShouldIgnore(relPath) {
			if verbose {
				fmt.Fprintf(os.Stderr, "[DEBUG] Skipping file: %s\n", relPath)
			}
			continue
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "[DEBUG] File: %s\n", relPath)
		}

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
					ifaceMethods[m] = methodInfo{
						ifaceName: tspec.Name.Name,
						iface:     ifaceType,
						method:    m,
						used:      false,
					}
				}
			}
		}
	}

	return ifaceMethods
}

// methodAnalyzer handles analysis of method usage in AST
type methodAnalyzer struct {
	pass         *analysis.Pass
	ifaceMethods map[*types.Func]methodInfo
	usedMethods  map[*types.Func]bool
}

// newMethodAnalyzer creates a new method analyzer
func newMethodAnalyzer(pass *analysis.Pass, ifaceMethods map[*types.Func]methodInfo) *methodAnalyzer {
	return &methodAnalyzer{
		pass:         pass,
		ifaceMethods: ifaceMethods,
		usedMethods:  make(map[*types.Func]bool),
	}
}

// analyzeUsedMethods traverses AST and marks used methods
func analyzeUsedMethods(pass *analysis.Pass, ifaceMethods map[*types.Func]methodInfo) map[*types.Func]bool {
	methodAnalyzer := newMethodAnalyzer(pass, ifaceMethods)
	return methodAnalyzer.analyze()
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
	for ifaceMethod, info := range ma.ifaceMethods {
		if ma.usedMethods[ifaceMethod] {
			continue
		}

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
	for ifaceMethod, info := range ma.ifaceMethods {
		if ma.usedMethods[ifaceMethod] {
			continue
		}

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

func Run() {
	singlechecker.Main(a)
}
