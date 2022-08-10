package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kffl/speedbump/lib"
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

	s, err := lib.NewSpeedbump(cfg)

	if err != nil {
		exitWithError(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool)

	go func() {
		<-sigs
		// signal was caught for the first time
		// stop the speedbump instance
		go func() {
			s.Stop()
			done <- true
		}()
		<-sigs
		// signal was caught for the second time
		// force the process to exit
		os.Exit(1)
	}()

	err = s.Start()

	if err != nil {
		exitWithError(err)
	}

	<-done
}
