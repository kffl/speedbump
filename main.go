package main

import (
	"fmt"
	"os"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func main() {
	cfg, err := parseArgs(os.Args[1:])

	if err != nil {
		exitWithError(err)
	}

	s, err := NewSpeedbump(cfg)

	if err != nil {
		exitWithError(err)
	}

	err = s.Start()

	if err != nil {
		exitWithError(err)
	}
}
