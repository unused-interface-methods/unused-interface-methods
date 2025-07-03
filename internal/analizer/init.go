package analizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unused-interface-methods/unused-interface-methods/internal/config"
)

var (
	verbose  bool
	basePath string
	cfg      *config.Config
)

func init() {
	var err error
	val := os.Getenv("UNUSED_INTERFACE_METHODS_VERBOSE")
	if val == "1" || val == "true" {
		verbose = true
	}
	basePath, err = extractBasePath(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting base path: %v\n", err)
		os.Exit(1)
	}
	cfg, err = config.LoadConfig("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
}

func extractBasePath(args []string) (string, error) {
	result := "."
	if len(args) > 0 {
		result = args[0]
		result = strings.TrimSuffix(result, "/...")
		result = strings.TrimPrefix(result, "./")
	}
	return filepath.Abs(result)
}
