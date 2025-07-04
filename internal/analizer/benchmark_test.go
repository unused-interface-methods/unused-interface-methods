package analizer

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func BenchmarkCollectInterfaceMethods(b *testing.B) {
	code := `
package test

type Interface1 interface {
	Method1() string
	Method2(int) error
	Method3(string, bool) (int, error)
}

type Interface2 interface {
	Method4() bool
	Method5([]byte) string
}

type Interface3 interface {
	Method6() interface{}
	Method7(interface{}) error
}
`

	for i := 0; i < b.N; i++ {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", code, 0)

		info := &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		}

		pkg, _ := new(types.Config).Check("test", fset, []*ast.File{file}, info)

		pass := &analysis.Pass{
			Fset:      fset,
			Files:     []*ast.File{file},
			Pkg:       pkg,
			TypesInfo: info,
		}

		collectInterfaceMethods(pass)
	}
}

func BenchmarkAnalyzeUsedMethods(b *testing.B) {
	code := `
package test

type Interface1 interface {
	Method1() string
	Method2(int) error
}

type Impl struct{}

func (i *Impl) Method1() string { return "" }
func (i *Impl) Method2(n int) error { return nil }

func useInterface() {
	var i Interface1 = &Impl{}
	i.Method1()
	i.Method2(42)
}
`

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, 0)

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	pkg, _ := new(types.Config).Check("test", fset, []*ast.File{file}, info)

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		ResultOf:  make(map[*analysis.Analyzer]interface{}),
	}

	// Add inspector result
	pass.ResultOf[inspect.Analyzer] = inspector.New([]*ast.File{file})

	ifaceMethods := collectInterfaceMethods(pass)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzeUsedMethods(pass, ifaceMethods)
	}
}

func BenchmarkMarkMatchingMethods(b *testing.B) {
	// Create a large number of interface methods to test performance
	ifaceMethods := make(map[*types.Func]methodInfo, 100)

	// Create a mock interface type for testing
	mockIface := types.NewInterfaceType(nil, nil)

	for i := 0; i < 100; i++ {
		sig := types.NewSignature(nil, nil, nil, false)
		method := types.NewFunc(token.NoPos, nil, fmt.Sprintf("Method%d", i), sig)
		ifaceMethods[method] = methodInfo{
			ifaceName: fmt.Sprintf("Interface%d", i/10),
			iface:     mockIface,
			method:    method,
		}
	}

	ma := &methodAnalyzer{
		ifaceMethods:  ifaceMethods,
		usedMethods:   make(map[*types.Func]bool),
		methodsByName: make(map[string][]*types.Func),
	}

	// Build method name index
	for method := range ifaceMethods {
		name := method.Name()
		ma.methodsByName[name] = append(ma.methodsByName[name], method)
	}

	// Create a test method to match
	testMethod := types.NewFunc(token.NoPos, nil, "Method50", types.NewSignature(nil, nil, nil, false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ma.usedMethods = make(map[*types.Func]bool) // Reset for each iteration
		ma.markMatchingMethods(testMethod, nil)
	}
}

// benchmarkTestdataFile runs benchmark on a specific testdata file
func benchmarkTestdataFile(b *testing.B, filename string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		b.Fatal(err)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	// Configure type checker with standard importer
	conf := &types.Config{
		Importer: importer.Default(),
	}

	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	if err != nil {
		b.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		ResultOf:  make(map[*analysis.Analyzer]interface{}),
	}

	// Add inspector result
	pass.ResultOf[inspect.Analyzer] = inspector.New([]*ast.File{file})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ifaceMethods := collectInterfaceMethods(pass)
		analyzeUsedMethods(pass, ifaceMethods)
	}
}

func BenchmarkTestdataInterfaces(b *testing.B) {
	benchmarkTestdataFile(b, "testdata/src/test/interfaces.go")
}

func BenchmarkTestdataReflection(b *testing.B) {
	benchmarkTestdataFile(b, "testdata/src/test/reflection.go")
}

func BenchmarkTestdataGenerics(b *testing.B) {
	benchmarkTestdataFile(b, "testdata/src/test/generics.go")
}

func BenchmarkTestdataAll(b *testing.B) {
	files := []string{
		"testdata/src/test/interfaces.go",
		"testdata/src/test/reflection.go",
		"testdata/src/test/generics.go",
	}

	fset := token.NewFileSet()
	var astFiles []*ast.File

	for _, filename := range files {
		file, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			b.Fatal(err)
		}
		astFiles = append(astFiles, file)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	// Configure type checker with standard importer
	conf := &types.Config{
		Importer: importer.Default(),
	}

	pkg, err := conf.Check("test", fset, astFiles, info)
	if err != nil {
		b.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     astFiles,
		Pkg:       pkg,
		TypesInfo: info,
		ResultOf:  make(map[*analysis.Analyzer]interface{}),
	}

	// Add inspector result
	pass.ResultOf[inspect.Analyzer] = inspector.New(astFiles)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ifaceMethods := collectInterfaceMethods(pass)
		analyzeUsedMethods(pass, ifaceMethods)
	}
}
