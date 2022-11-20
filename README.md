# speedbump - TCP proxy with variable latency

<div align="center">
  <img alt="speedbump logo" src="https://github.com/kffl/speedbump/raw/HEAD/assets/speedbump.gif" width="480" height="auto"/>
</div>
Speedbump is a TCP proxy written in Go which allows for simulating variable network latency.

[![CI Workflow](https://github.com/kffl/speedbump/workflows/CI/badge.svg)](https://github.com/kffl/speedbump/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/kffl/speedbump)](https://goreportcard.com/report/github.com/kffl/speedbump) [![Docker Pulls](https://img.shields.io/docker/pulls/kffl/speedbump)](https://hub.docker.com/r/kffl/speedbump) [![Docker Image Version](https://img.shields.io/docker/v/kffl/speedbump)](https://hub.docker.com/r/kffl/speedbump) [![GoDoc](https://godoc.org/github.com/kffl/speedbump/lib?status.svg)](https://godoc.org/github.com/kffl/speedbump/lib)

## Usage

### Installation

The easiest way to install speedbump is to download pre-built binaries for your platform that are automatically attached to each [release](https://github.com/kffl/speedbump/releases/) under _Assets_. If you wish to build speedbump from source, clone this repository and run `go build`. Alternatively, you can run speedbump as a container using the [kffl/speedbump](https://hub.docker.com/r/kffl/speedbump) image.

### Basic usage examples

Spawn a new instance listening on port 2000 that proxies TCP traffic to localhost:80 with a base latency of 100ms and sine wave amplitude of 100ms (resulting in maximum added latency being 200ms and minimum being 0), period of which is 1 minute:

```
speedbump --latency=100ms --sine-amplitude=100ms --sine-period=1m --port=2000 localhost:80
```

or when running speedbump using the [kffl/speedbump](https://hub.docker.com/r/kffl/speedbump) container image:

```
docker run --net=host kffl/speedbump:latest --latency=100ms --sine-amplitude=100ms \
           --sine-period=1m --port=2000 localhost:80
```

Spawn a new instance with a base latency of 300ms and a sawtooth wave latency summand with amplitude of 200ms and period of 2 minutes (visualized by the graph below):

```
speedbump --latency=300ms --saw-amplitude=200ms --saw-period=2m --port=2000 localhost:80
```

<div align="center">
  <img alt="speedbump sawtooth wave graph" src="https://github.com/kffl/speedbump/raw/HEAD/assets/sawtooth.svg" width="800" height="auto"/>
</div>

### Combining latency summands

It is possible to run speedbump with multiple latency summands at once:

<div align="center">
  <img alt="speedbump sawtooth + sine graph" src="https://github.com/kffl/speedbump/raw/HEAD/assets/combined.svg" width="800" height="auto"/>
</div>

## CLI Arguments Reference:

Output of `speedbump --help`:

```
usage: speedbump [<flags>] <destination>

TCP proxy for simulating variable network latency.

Flags:
  --help                  Show context-sensitive help (also try --help-long and
                          --help-man).
  --host=""               IP or hostname to listen on. Speedbump will bind to
                          all available network interfaces if unspecified.
  --port=8000             Port number to listen on.
  --buffer=64KB           Size of the buffer used for TCP reads.
  --queue-size=1024       Size of the delay queue storing read buffers.
  --latency=5ms           Base latency added to proxied traffic.
  --log-level=INFO        Log level. Possible values: DEBUG, TRACE, INFO, WARN,
                          ERROR.
  --sine-amplitude=0      Amplitude of the latency sine wave.
  --sine-period=0         Period of the latency sine wave.
  --saw-amplitude=0       Amplitude of the latency sawtooth wave.
  --saw-period=0          Period of the latency sawtooth wave.
  --square-amplitude=0    Amplitude of the latency square wave.
  --square-period=0       Period of the latency square wave.
  --triangle-amplitude=0  Amplitude of the latency triangle wave.
  --triangle-period=0     Period of the latency triangle wave.
  --version               Show application version.

Args:
  <destination>  TCP proxy destination in host:post format.
```

## Using speedbump as a library

Speedbump can be used as a Go library via its `lib` package. Check `lib` [README](lib/README.md) for additional information.

## License

Copyright Paweł Kuffel 2022, licensed under Apache 2.0 License.

Speedbump logo contains the Go Gopher mascot which was originally designed by Renee French (http://reneefrench.blogspot.com/) and licensed under Creative Commons 3.0 Attributions license.
