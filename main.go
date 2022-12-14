package main

import (
	"os"

	"github.com/fatih/color"

	"shipyard/cmd"
)

func main() {
	// Handle a panic.
	defer func() {
		if err := recover(); err != nil {
			red := color.New(color.FgHiRed)
			red.Fprintf(os.Stderr, "Runtime error: %v\n", err)
			os.Exit(1)
		}

	}()
	cmd.Execute()
}
