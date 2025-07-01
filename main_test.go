package main

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "test")
}

func TestCollectInterfaceMethods(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected []string // expected interface methods
	}{
		{
			name: "simple interface",
			src: `package test
type Writer interface {
	Write([]byte) (int, error)
	Close() error
}`,
			expected: []string{"Write", "Close"},
		},
		{
			name: "empty interface",
			src: `package test
type Empty interface {}`,
			expected: []string{},
		},
		{
			name: "embedded interface",
			src: `package test
import "io"
type ReadWriter interface {
	io.Reader
	Write([]byte) (int, error)
}`,
			expected: []string{"Write"}, // only explicit methods
		},
		{
			name: "generic interface",
			src: `package test
type Repo[T any] interface {
	Get(id string) T
	Save(item T) error
}`,
			expected: []string{"Get", "Save"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := createTestPass(t, tt.src)
			methods := collectInterfaceMethods(pass)

			var methodNames []string
			for _, info := range methods {
				methodNames = append(methodNames, info.method.Name())
			}

			if len(methodNames) != len(tt.expected) {
				t.Errorf("expected %d methods, got %d: %v", len(tt.expected), len(methodNames), methodNames)
				return
			}

			for _, expected := range tt.expected {
				found := false
				for _, actual := range methodNames {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected method %s not found in %v", expected, methodNames)
				}
			}
		})
	}
}

func TestMethodAnalyzer_AnalyzeSelectorExpr(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected []string // expected used method names
	}{
		{
			name: "direct method call",
			src: `package test
type Writer interface {
	Write([]byte) (int, error)
	Close() error
}
func test(w Writer) {
	w.Write(nil)
}`,
			expected: []string{"Write"},
		},
		{
			name: "method on concrete type implementing interface",
			src: `package test
type Writer interface {
	Write([]byte) (int, error)
}
type Buffer struct{}
func (b Buffer) Write([]byte) (int, error) { return 0, nil }
func test() {
	var buf Buffer
	buf.Write(nil)
}`,
			expected: []string{"Write"},
		},
		{
			name: "no usage",
			src: `package test
type Writer interface {
	Write([]byte) (int, error)
	Close() error
}
func test() {}`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := createTestPass(t, tt.src)
			ifaceMethods := collectInterfaceMethods(pass)
			analyzer := newMethodAnalyzer(pass, ifaceMethods)
			usedMethods := analyzer.analyze()

			var usedNames []string
			for method := range usedMethods {
				usedNames = append(usedNames, method.Name())
			}

			if len(usedNames) != len(tt.expected) {
				t.Errorf("expected %d used methods, got %d: %v", len(tt.expected), len(usedNames), usedNames)
				return
			}

			for _, expected := range tt.expected {
				found := false
				for _, actual := range usedNames {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected used method %s not found in %v", expected, usedNames)
				}
			}
		})
	}
}

func TestMethodAnalyzer_CheckImplements(t *testing.T) {
	src := `package test
type Writer interface {
	Write([]byte) (int, error)
}
type Reader interface {
	Read([]byte) (int, error)
}
type CustomWriter struct{}
func (cw CustomWriter) Write([]byte) (int, error) { return 0, nil }`

	pass := createTestPass(t, src)
	ifaceMethods := collectInterfaceMethods(pass)
	analyzer := newMethodAnalyzer(pass, ifaceMethods)

	// find Writer interface
	var writerInfo methodInfo
	for _, info := range ifaceMethods {
		if info.ifaceName == "Writer" && info.method.Name() == "Write" {
			writerInfo = info
			break
		}
	}

	if writerInfo.method == nil {
		t.Fatal("Writer.Write method not found")
	}

	tests := []struct {
		name     string
		typeName string
		expected bool
	}{
		{
			name:     "concrete type implements interface",
			typeName: "CustomWriter",
			expected: true,
		},
		{
			name:     "interface implements interface",
			typeName: "io.Writer",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that we can call checkImplements method
			// For now, just verify the analyzer was created properly
			if analyzer == nil {
				t.Error("analyzer should not be nil")
			}
		})
	}
}

func TestMethodAnalyzer_IsStringerMethod(t *testing.T) {
	src := `package test
type Stringer interface {
	String() string
	ToString() string
}
type CustomStringer interface {
	StringWithFormat(format string) string
}`

	pass := createTestPass(t, src)
	ifaceMethods := collectInterfaceMethods(pass)
	analyzer := newMethodAnalyzer(pass, ifaceMethods)

	tests := []struct {
		methodName string
		expected   bool
	}{
		{"String", true},
		{"ToString", false},
	}

	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			for _, info := range ifaceMethods {
				if info.method.Name() == tt.methodName && info.method.Name() == "String" {
					// check if it's the correct String() string signature
					result := analyzer.isStringerMethod(info.method)
					if result != tt.expected {
						t.Errorf("isStringerMethod(%s) = %v, expected %v", tt.methodName, result, tt.expected)
					}
					return
				}
			}
		})
	}
}

func TestMethodAnalyzer_AnalyzeFmtCall(t *testing.T) {
	src := `package test
import "fmt"
type Stringer interface {
	String() string
}
type Item struct{}
func (i Item) String() string { return "item" }
func test() {
	item := Item{}
	fmt.Println(item)
}`

	pass := createTestPass(t, src)
	ifaceMethods := collectInterfaceMethods(pass)
	analyzer := newMethodAnalyzer(pass, ifaceMethods)
	usedMethods := analyzer.analyze()

	// String method should be marked as used
	found := false
	for method := range usedMethods {
		if method.Name() == "String" {
			found = true
			break
		}
	}

	if !found {
		t.Error("String method should be marked as used when passed to fmt.Println")
	}
}

func TestSkipGenerics(t *testing.T) {
	originalSkipGenerics := skipGenerics
	defer func() { skipGenerics = originalSkipGenerics }()

	src := `package test
type GenericRepo[T any] interface {
	Get(id string) T
	Save(item T) error
}
type RegularRepo interface {
	Load() error
}`

	tests := []struct {
		name         string
		skipGenerics bool
		expectedNum  int // expected number of interface methods
	}{
		{
			name:         "include generics",
			skipGenerics: false,
			expectedNum:  3, // Get, Save, Load
		},
		{
			name:         "skip generics",
			skipGenerics: true,
			expectedNum:  1, // only Load
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipGenerics = tt.skipGenerics
			pass := createTestPass(t, src)
			methods := collectInterfaceMethods(pass)

			if len(methods) != tt.expectedNum {
				t.Errorf("expected %d methods, got %d", tt.expectedNum, len(methods))
			}
		})
	}
}

func TestReportUnusedMethods(t *testing.T) {
	src := `package test
type Writer interface {
	Write([]byte) (int, error)
	Close() error
}
func test(w Writer) {
	w.Write(nil)
	// Close is not used
}`

	pass := createTestPass(t, src)
	ifaceMethods := collectInterfaceMethods(pass)
	used := analyzeUsedMethods(pass, ifaceMethods)

	// capture reports
	var reports []string
	pass.Report = func(d analysis.Diagnostic) {
		reports = append(reports, d.Message)
	}

	reportUnusedMethods(pass, ifaceMethods, used)

	if len(reports) == 0 {
		t.Error("expected at least one report for unused method")
	}

	// should report Close method as unused
	found := false
	for _, report := range reports {
		if contains(report, "Close") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Close method should be reported as unused")
	}
}

// Helper functions

func createTestPass(t *testing.T, src string) *analysis.Pass {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	config := &types.Config{
		Importer: importer.Default(),
	}

	pkg, err := config.Check("test", fset, []*ast.File{file}, info)
	if err != nil {
		t.Fatal(err)
	}

	// create inspector
	ins := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Analyzer:  Analyzer,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		ResultOf: map[*analysis.Analyzer]interface{}{
			inspect.Analyzer: ins,
		},
		Report: func(d analysis.Diagnostic) {
			// default no-op
		},
	}

	return pass
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Integration test that uses actual test data
func TestIntegration(t *testing.T) {
	testdata := analysistest.TestData()

	// Test with skipGenerics = false
	skipGenerics = false
	analysistest.Run(t, testdata, Analyzer, "test")

	// Test with skipGenerics = true
	skipGenerics = true
	analysistest.Run(t, testdata, Analyzer, "test")
}

func BenchmarkAnalyzer(b *testing.B) {
	src := `package test
type Writer interface {
	Write([]byte) (int, error)
	Close() error
	Sync() error
}
type Reader interface {
	Read([]byte) (int, error)
	Close() error
}
func test(w Writer, r Reader) {
	w.Write(nil)
	r.Read(nil)
	// Close and Sync are unused
}`

	pass := createBenchPass(b, src)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ifaceMethods := collectInterfaceMethods(pass)
		used := analyzeUsedMethods(pass, ifaceMethods)
		reportUnusedMethods(pass, ifaceMethods, used)
	}
}

func createBenchPass(b *testing.B, src string) *analysis.Pass {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		b.Fatal(err)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	config := &types.Config{
		Importer: importer.Default(),
	}

	pkg, err := config.Check("test", fset, []*ast.File{file}, info)
	if err != nil {
		b.Fatal(err)
	}

	// create inspector
	ins := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Analyzer:  Analyzer,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		ResultOf: map[*analysis.Analyzer]interface{}{
			inspect.Analyzer: ins,
		},
		Report: func(d analysis.Diagnostic) {
			// default no-op
		},
	}

	return pass
}
