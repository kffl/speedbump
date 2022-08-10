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
		QueueSize:  2048,
		Latency: &speedbump.LatencyCfg{
			Base:          time.Millisecond * 100,
			SineAmplitude: time.Millisecond * 50,
			SinePeriod:    time.Minute,
		},
		LogLevel: "TRACE",
	}

	s, err := speedbump.NewSpeedbump(&cfg)

	if err != nil {
		// handle creation error
		return
	}

	// Start() will unblock as soon as the proxy is started
	// or return an error if there is a startup error
	err = s.Start()

	if err != nil {
		// handle startup error
		return
	}

	// let's stop the proxy after 5 mins
	time.Sleep(time.Minute * 5)

	s.Stop()

	// DONE
}

```

## `v1` Upgrade guide

In an effort to make the `lib` package easier to work with when used as a dependency for Go tests, the following changes were made to its API in the `v1` release:

- `Start()` is no longer blocking. It will either unblock as soon as the proxy starts listening or return an error if proxy startup fails;
- `Stop()` waits for all proxy connections to close before returning;
- a field name typo `sawAmplitute` was fixed in `LatencyCfg` struct (renamed to `SawAmplitude`).