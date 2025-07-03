package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "test")
}

func BenchmarkAnalyzer(b *testing.B) {
	testdata := analysistest.TestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analysistest.Run(&testing.T{}, testdata, Analyzer, "test")
	}
}

// sequential version for comparison
func collectInterfaceMethodsSequential(pass *analysis.Pass) map[*types.Func]methodInfo {
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

// create analyzer with sequential version
var SequentialAnalyzer = &analysis.Analyzer{
	Name:     "unused_interface_methods_sequential",
	Doc:      "Sequential version for benchmarking",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run: func(pass *analysis.Pass) (interface{}, error) {
		ifaceMethods := collectInterfaceMethodsSequential(pass)
		used := analyzeUsedMethods(pass, ifaceMethods)
		reportUnusedMethods(pass, ifaceMethods, used)
		return nil, nil
	},
}

func BenchmarkAnalyzerSequential(b *testing.B) {
	testdata := analysistest.TestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analysistest.Run(&testing.T{}, testdata, SequentialAnalyzer, "test")
	}
}
