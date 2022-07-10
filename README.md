# speedbump - TCP proxy for simulating variable network latency

## Usage

Spawn a new instance listening on port 2000 that proxies TCP traffic to localhost:80 with a base latency of 100ms and sine wave amplitude of 100ms (resulting in maximum added latency being 200ms and minimum being 0), period of which is 1m:

```
speedbump --latency=100ms --sine-amplitude=100ms --sine-period=1m --port=2000 localhost:80
```

## CLI Arguments Reference:

Output of `speedbump --help`:

```
usage: speedbump [<flags>] <destination>

TCP proxy for simulating variable network latency.

Flags:
  --help              Show context-sensitive help (also try --help-long and --help-man).
  --port=8000         Port number to listen on.
  --buffer=64KB       Size of the buffer used for TCP reads.
  --latency=5ms       Base latency added to proxied traffic.
  --sine-amplitude=0  Amplitude of the latency sine wave.
  --sine-period=0     Period of the latency sine wave.
  --saw-amplitude=0   Amplitude of the latency sawtooth wave.
  --saw-period=0      Period of the latency sawtooth wave.
  --version           Show application version.

Args:
  <destination>  TCP proxy destination in host:post format.
```