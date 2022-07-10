package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func parseArgs(args []string) (*SpeedbumpCfg, error) {
	var app = kingpin.New("speedbump", "TCP proxy for simulating variable network latency.")

	var (
		port       = app.Flag("port", "Port number to listen on.").Default("8000").Int()
		bufferSize = app.Flag("buffer", "Size of the buffer used for TCP reads.").
				Default("64KB").
				Bytes()
		latency = app.Flag("latency", "Base latency added to proxied traffic.").
			Default("5ms").
			Duration()
		sineAmplitude = app.Flag("sine-amplitude", "Amplitude of the latency sine wave.").
				PlaceHolder("0").
				Duration()
		sinePeriod = app.Flag("sine-period", "Period of the latency sine wave.").
				PlaceHolder("0").
				Duration()
		sawAmplitute = app.Flag("saw-amplitude", "Amplitude of the latency sawtooth wave.").
				PlaceHolder("0").
				Duration()
		sawPeriod = app.Flag("saw-period", "Period of the latency sawtooth wave.").
				PlaceHolder("0").
				Duration()
		destAddr = app.Arg("destination", "TCP proxy destination in host:post format.").
				Required().
				String()
	)

	app.Version("0.1.0-rc2")
	_, err := app.Parse(args)

	if err != nil {
		return nil, err
	}

	var cfg = SpeedbumpCfg{
		Port:       *port,
		DestAddr:   *destAddr,
		BufferSize: int(*bufferSize),
		Latency: &LatencyCfg{
			base:          *latency,
			sineAmplitude: *sineAmplitude,
			sinePeriod:    *sinePeriod,
			sawAmplitute:  *sawAmplitute,
			sawPeriod:     *sawPeriod,
		},
	}

	return &cfg, err
}
