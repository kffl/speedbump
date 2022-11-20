package main

import (
	"github.com/kffl/speedbump/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

func parseArgs(args []string) (*lib.SpeedbumpCfg, error) {
	var app = kingpin.New("speedbump", "TCP proxy for simulating variable network latency.")

	var (
		host = app.Flag("host", "IP or hostname to listen on. Speedbump will bind to all network interfaces if unspecified.").
			Default("").
			String()
		port       = app.Flag("port", "Port number to listen on.").Default("8000").Int()
		bufferSize = app.Flag("buffer", "Size of the buffer used for TCP reads.").
				Default("64KB").
				Bytes()
		queueSize = app.Flag("queue-size", "Size of the delay queue storing read buffers.").
				Default("1024").
				Int()
		latency = app.Flag("latency", "Base latency added to proxied traffic.").
			Default("5ms").
			Duration()
		logLevel = app.Flag("log-level", "Log level. Possible values: DEBUG, TRACE, INFO, WARN, ERROR.").
				Default("INFO").
				Enum("DEBUG", "TRACE", "INFO", "WARN", "ERROR")
		sineAmplitude = app.Flag("sine-amplitude", "Amplitude of the latency sine wave.").
				PlaceHolder("0").
				Duration()
		sinePeriod = app.Flag("sine-period", "Period of the latency sine wave.").
				PlaceHolder("0").
				Duration()
		sawAmplitude = app.Flag("saw-amplitude", "Amplitude of the latency sawtooth wave.").
				PlaceHolder("0").
				Duration()
		sawPeriod = app.Flag("saw-period", "Period of the latency sawtooth wave.").
				PlaceHolder("0").
				Duration()
		squareAmplitude = app.Flag("square-amplitude", "Amplitude of the latency square wave.").
				PlaceHolder("0").
				Duration()
		squarePeriod = app.Flag("square-period", "Period of the latency square wave.").
				PlaceHolder("0").
				Duration()
		triangleAmplitude = app.Flag("triangle-amplitude", "Amplitude of the latency triangle wave.").
					PlaceHolder("0").
					Duration()
		trianglePeriod = app.Flag("triangle-period", "Period of the latency triangle wave.").
				PlaceHolder("0").
				Duration()
		destAddr = app.Arg("destination", "TCP proxy destination in host:post format.").
				Required().
				String()
	)

	app.Version("1.1.0")
	_, err := app.Parse(args)

	if err != nil {
		return nil, err
	}

	var cfg = lib.SpeedbumpCfg{
		Host:       *host,
		Port:       *port,
		DestAddr:   *destAddr,
		BufferSize: int(*bufferSize),
		QueueSize:  *queueSize,
		Latency: &lib.LatencyCfg{
			Base:              *latency,
			SineAmplitude:     *sineAmplitude,
			SinePeriod:        *sinePeriod,
			SawAmplitude:      *sawAmplitude,
			SawPeriod:         *sawPeriod,
			SquareAmplitude:   *squareAmplitude,
			SquarePeriod:      *squarePeriod,
			TriangleAmplitude: *triangleAmplitude,
			TrianglePeriod:    *trianglePeriod,
		},
		LogLevel: *logLevel,
	}

	return &cfg, err
}
