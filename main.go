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

var skipGenerics bool

func init() {
	Analyzer.Flags.BoolVar(&skipGenerics, "skipGenerics", false, "skip interfaces with type parameters (generics)")
}

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

				// Опционально пропускаем дженерик-интерфейсы (с параметрами типа)
				if skipGenerics {
					if named.TypeParams() != nil && named.TypeParams().Len() > 0 {
						continue
					}
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

	// Шаг B: отметить использованные методы
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.SelectorExpr)(nil),
		(*ast.CallExpr)(nil),
	}

	usedMethods := make(map[*types.Func]bool)

	ins.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			// Обработка прямых вызовов методов
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

				if ifaceMethod == calledMethod {
					usedMethods[ifaceMethod] = true
					continue
				}

				if calledMethod.Name() != ifaceMethod.Name() ||
					calledMethod.Type().(*types.Signature).String() != ifaceMethod.Type().(*types.Signature).String() {
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

				// --- Generic fallback: match by name and signature for instantiated generic interfaces ---
				if calledMethod.Name() == ifaceMethod.Name() &&
					calledMethod.Type().(*types.Signature).String() == ifaceMethod.Type().(*types.Signature).String() {
					usedMethods[ifaceMethod] = true
				}
			}

		case *ast.CallExpr:
			// Обработка неявных вызовов для fmt.Stringer
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
					// Проверка на сигнатуру `String() string`
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

	// Применяем результаты анализа
	for ifaceMethod := range usedMethods {
		if info, exists := ifaceMethods[ifaceMethod]; exists {
			info.used = true
			ifaceMethods[ifaceMethod] = info
		}
	}

	// Шаг C: собрать, отсортировать и вывести отчет по неиспользуемым методам
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
