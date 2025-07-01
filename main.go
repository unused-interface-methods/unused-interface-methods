package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer реализует плагины для поиска неиспользуемых методов интерфейсов.
var Analyzer = &analysis.Analyzer{
	Name:     "unusedint",
	Doc:      "находит методы интерфейсов, объявленные, но не используемые в коде",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

type methodInfo struct {
	ifaceName string           // имя интерфейса
	iface     *types.Interface // объект интерфейса
	method    *types.Func      // объект метода
	used      bool             // флаг использования
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Шаг A: собрать информацию об интерфейсных методах
	// map[method] -> list of interfaces that have this method
	ifaceMethods := map[*types.Func][]methodInfo{}

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
				ifaceType, ok := named.Underlying().(*types.Interface)
				if !ok {
					continue
				}

				for i := 0; i < ifaceType.NumMethods(); i++ {
					m := ifaceType.Method(i)
					if m == nil {
						continue
					}
					info := methodInfo{
						ifaceName: tspec.Name.Name,
						iface:     ifaceType,
						method:    m,
						used:      false,
					}
					ifaceMethods[m] = append(ifaceMethods[m], info)
				}
			}
		}
	}

	// Шаг B: отметить использованные методы
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
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

			if infos, ok := ifaceMethods[calledMethod]; ok {
				for i := range infos {
					info := &infos[i]
					if !types.IsInterface(recv) {
						if types.Implements(recv, info.iface) {
							infos[i].used = true
						}
					} else {
						if recvIface, ok := recv.Underlying().(*types.Interface); ok {
							if types.Implements(info.iface, recvIface) {
								infos[i].used = true
							}
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

				for ifaceMethod, infos := range ifaceMethods {
					sig, ok := ifaceMethod.Type().(*types.Signature)
					if !ok || ifaceMethod.Name() != "String" || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
						continue
					}
					if basic, ok := sig.Results().At(0).Type().(*types.Basic); !ok || basic.Kind() != types.String {
						continue
					}

					for i := range infos {
						info := &infos[i]
						if types.Implements(argType, info.iface) {
							infos[i].used = true
						}
					}
				}
			}
		}
	})

	// Шаг C: собрать, отсортировать и вывести отчет по неиспользуемым методам
	var unused []methodInfo
	for _, infos := range ifaceMethods {
		for _, info := range infos {
			if !info.used {
				unused = append(unused, info)
			}
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
		pass.Reportf(info.method.Pos(),
			"метод %q интерфейса %q объявлен, но не используется",
			info.method.Name(),
			info.ifaceName,
		)
	}
	return nil, nil
}

func main() {
	multichecker.Main(Analyzer)
}
