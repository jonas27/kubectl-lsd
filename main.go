package main

import (
	"fmt"
	"os"

	"github.com/jonas27/kubectl-lsd/cmd"
)

const (
	exitCodeErr = 1
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "LSD stopped with error: %v\n", err)
		os.Exit(exitCodeErr)
	}
}
