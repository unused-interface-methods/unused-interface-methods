package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

var skipGenerics bool

func init() {
	Analyzer.Flags.BoolVar(&skipGenerics, "skipGenerics", false, "skip interfaces with type parameters (generics)")
}

// Analyzer implements plugins for finding unused interface methods.
var Analyzer = &analysis.Analyzer{
	Name:     "unusedintf",
	Doc:      "finds interface methods that are declared but not used in the code",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

type methodInfo struct {
	ifaceName string           // имя интерфейса
	iface     *types.Interface // объект интерфейса
	method    *types.Func      // объект метода
	used      bool             // флаг использования
}

// collectInterfaceMethods collects all explicit interface methods in the package.
func collectInterfaceMethods(pass *analysis.Pass) map[*types.Func]methodInfo {
	ifaceMethods := map[*types.Func]methodInfo{}

	for _, file := range pass.Files {
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

				// Skip generic interfaces when necessary.
				if skipGenerics && named.TypeParams() != nil && named.TypeParams().Len() > 0 {
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

// analyzeUsedMethods traverses AST and marks used methods.
func analyzeUsedMethods(pass *analysis.Pass, ifaceMethods map[*types.Func]methodInfo) map[*types.Func]bool {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	usedMethods := make(map[*types.Func]bool)

	nodeFilter := []ast.Node{
		(*ast.SelectorExpr)(nil),
		(*ast.CallExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			sel := pass.TypesInfo.Selections[node]
			if sel == nil || (sel.Kind() != types.MethodVal && sel.Kind() != types.MethodExpr) {
				return
			}

			calledMethod := sel.Obj().(*types.Func)
			recv := sel.Recv()

			for ifaceMethod, info := range ifaceMethods {
				if usedMethods[ifaceMethod] {
					continue
				}

				if calledMethod == ifaceMethod {
					usedMethods[ifaceMethod] = true
					continue
				}

				if calledMethod.Name() != ifaceMethod.Name() {
					continue
				}

				implemented := types.Implements(recv, info.iface)
				if !implemented {
					if recvIface, ok := recv.Underlying().(*types.Interface); ok {
						if types.Implements(info.iface, recvIface) {
							implemented = true
						}
					}
				}

				if implemented {
					usedMethods[ifaceMethod] = true
					continue
				}

				// Generic heuristic: имя совпало, сигнатуры совместимы, а получатель является инстанциацией iface
				if namedRecv, ok := recv.(*types.Named); ok {
					if origin := namedRecv.Origin(); origin != nil {
						if ifaceOrig, ok2 := origin.Underlying().(*types.Interface); ok2 && types.Identical(ifaceOrig, info.iface) {
							usedMethods[ifaceMethod] = true
						}
					}
				}
			}

		case *ast.CallExpr:
			var ident *ast.Ident
			switch fun := node.Fun.(type) {
			case *ast.Ident:
				ident = fun
			case *ast.SelectorExpr:
				ident = fun.Sel
			default:
				return
			}

			fn, ok := pass.TypesInfo.Uses[ident].(*types.Func)
			if !ok || fn.Pkg() == nil || fn.Pkg().Path() != "fmt" {
				return
			}

			for _, arg := range node.Args {
				argType := pass.TypesInfo.TypeOf(arg)
				if argType == nil {
					continue
				}

				for ifaceMethod, info := range ifaceMethods {
					if usedMethods[ifaceMethod] {
						continue
					}
					sig, ok := ifaceMethod.Type().(*types.Signature)
					if !ok || ifaceMethod.Name() != "String" || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
						continue
					}
					if basic, ok := sig.Results().At(0).Type().(*types.Basic); !ok || basic.Kind() != types.String {
						continue
					}

					if types.Implements(argType, info.iface) {
						usedMethods[ifaceMethod] = true
					}
				}
			}
		}
	})

	return usedMethods
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
