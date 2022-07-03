package main

import (
	"math"
	"time"
)

type LatencyGenerator interface {
	generateLatency(when time.Time) time.Duration
}

type LatencyCfg struct {
	base          time.Duration
	sineAmplitude time.Duration
	sinePeriod    time.Duration
}

type simpleLatencyGenerator struct {
	start time.Time
	cfg   *LatencyCfg
}

func (g *simpleLatencyGenerator) generateLatency(when time.Time) time.Duration {

	return g.cfg.base + time.Duration(
		math.Sin(
			float64(when.Sub(g.start))/float64(g.cfg.sinePeriod)*math.Pi*2,
		)*float64(
			g.cfg.sineAmplitude,
		),
	)
}
