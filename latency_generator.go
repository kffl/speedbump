package main

import (
	"time"
)

type LatencyGenerator interface {
	generateLatency(time.Time) time.Duration
}

type LatencyCfg struct {
	base          time.Duration
	sineAmplitude time.Duration
	sinePeriod    time.Duration
	sawAmplitute  time.Duration
	sawPeriod     time.Duration
}

type latencySummand interface {
	getLatency(elapsed time.Duration) time.Duration
}

type simpleLatencyGenerator struct {
	start    time.Time
	summands []latencySummand
}

func newSimpleLatencyGenerator(start time.Time, cfg *LatencyCfg) simpleLatencyGenerator {
	summands := []latencySummand{baseLatencySummand{cfg.base}}
	if cfg.sineAmplitude > 0 && cfg.sinePeriod > 0 {
		summands = append(summands, sineLatencySummand{
			cfg.sineAmplitude,
			cfg.sinePeriod,
		})
	}
	if cfg.sawAmplitute > 0 && cfg.sawPeriod > 0 {
		summands = append(summands, sawtoothLatencySummand{
			cfg.sawAmplitute,
			cfg.sawPeriod,
		})
	}
	return simpleLatencyGenerator{
		start:    start,
		summands: summands,
	}
}

func (g simpleLatencyGenerator) generateLatency(when time.Time) time.Duration {
	var latency time.Duration = 0
	elapsed := when.Sub(g.start)
	for _, s := range g.summands {
		latency += s.getLatency(elapsed)
	}
	return latency
}
