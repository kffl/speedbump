package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		go s.Stop()
		<-sigs
		os.Exit(1)
	}()

	err = s.Start()

	if err != nil {
		exitWithError(err)
	}
}
