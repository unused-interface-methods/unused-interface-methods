package main

import (
	"flag"
	"fmt"
	"runtime/debug"

	"github.com/unused-interface-methods/unused-interface-methods/internal/analizer"
)

func main() {
	vFlag := flag.Bool("v", false, "show version")
	flag.Parse()
	if *vFlag {
			info, ok := debug.ReadBuildInfo()
			if ok {
					fmt.Println("Version:", info.Main.Version)
			} else {
					fmt.Println("Version: unknown")
			}
			return
	}
	analizer.Run()
}
