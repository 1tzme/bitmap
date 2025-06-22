package main

import (
	"fmt"
	"os"

	b "bitmap/internal/bmp"
	u "bitmap/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		u.PrintUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "header":
		b.HandleHeaderCommand()
	default:
		u.PrintUsage()
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
