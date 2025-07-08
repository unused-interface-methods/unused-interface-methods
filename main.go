package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/unused-interface-methods/unused-interface-methods/internal/analizer"
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			info, ok := debug.ReadBuildInfo()
			if ok && info.Main.Version != "" {
				fmt.Println("Version:", info.Main.Version)
			} else {
				fmt.Println("Version: unknown")
			}
			return
		}
	}
	analizer.Run()
}
