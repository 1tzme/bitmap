package main

import (
	"fmt"
	"os"
	u "bitmap/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		u.PrintUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	default:
		u.PrintUsage()
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}