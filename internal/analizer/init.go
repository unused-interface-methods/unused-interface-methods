package analizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unused-interface-methods/unused-interface-methods/internal/config"
)

var (
	basePath string
	cfg      *config.Config
)

func init() {
	var err error
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
