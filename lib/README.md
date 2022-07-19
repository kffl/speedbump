# speedbump `lib` package

This package allows for using speedbump as external library in Go code. It can be useful for adding programmatic delay to TCP connections while running load tests (i.e. between the SUT and a database).

Consult GoDoc for API reference:

[![GoDoc](https://godoc.org/github.com/kffl/speedbump/lib?status.svg)](https://godoc.org/github.com/kffl/speedbump/lib)

## Installation

```
go get github.com/kffl/speedbump/lib
```

## Example usage

```go
package main

import (
	"time"

	speedbump "github.com/kffl/speedbump/lib"
)

func main() {
	cfg := speedbump.SpeedbumpCfg{
		Port:       8000,
		DestAddr:   "localhost:80",
		BufferSize: 16384,
		Latency: &speedbump.LatencyCfg{
			Base:          time.Millisecond * 100,
			SineAmplitude: time.Millisecond * 50,
			SinePeriod:    time.Minute,
		},
		LogLevel: "INFO",
	}

	s, err := speedbump.NewSpeedbump(&cfg)

	if err != nil {
		// handle creation error
		return
	}

	go func() {
		// stop the proxy after 5 minutes
		time.Sleep(time.Second * 5)
		s.Stop()
	}()

	// Start() will either block until .Stop() is called
	// or return immedietely if there is a startup error
	err = s.Start()

	if err != nil {
		// handle startup error
		return
	}

	// DONE
}
```

