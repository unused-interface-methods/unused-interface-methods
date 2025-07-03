package analizer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	// It is not a bug, it is a specific analysistest. ))
	// [DEBUG] File: ../../../../../../../../../testdata/src/test/generics.go
	// [DEBUG] File: ../../../../../../../../../testdata/src/test/interfaces.go
	// [DEBUG] File: ../../../../../../../../../testdata/src/test/reflection.go
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, a, "test")
}
