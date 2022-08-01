package lib

import (
	"time"
)

type LatencyGenerator interface {
	generateLatency(time.Time) time.Duration
}

type LatencyCfg struct {
	Base            time.Duration
	SineAmplitude   time.Duration
	SinePeriod      time.Duration
	SawAmplitute    time.Duration
	SawPeriod       time.Duration
	SquareAmplitude time.Duration
	SquarePeriod    time.Duration
}

type latencySummand interface {
	getLatency(elapsed time.Duration) time.Duration
}

type simpleLatencyGenerator struct {
	start    time.Time
	summands []latencySummand
}

func newSimpleLatencyGenerator(start time.Time, cfg *LatencyCfg) simpleLatencyGenerator {
	summands := []latencySummand{baseLatencySummand{cfg.Base}}
	if cfg.SineAmplitude > 0 && cfg.SinePeriod > 0 {
		summands = append(summands, sineLatencySummand{
			cfg.SineAmplitude,
			cfg.SinePeriod,
		})
	}
	if cfg.SawAmplitute > 0 && cfg.SawPeriod > 0 {
		summands = append(summands, sawtoothLatencySummand{
			cfg.SawAmplitute,
			cfg.SawPeriod,
		})
	}
	if cfg.SquareAmplitude > 0 && cfg.SquarePeriod > 0 {
		summands = append(summands, squareLatencySummand{
			cfg.SquareAmplitude,
			cfg.SquarePeriod,
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
