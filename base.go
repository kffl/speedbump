package main

import "time"

type baseLatencySummand struct {
	latency time.Duration
}

func (b baseLatencySummand) getLatency(elapsed time.Duration) time.Duration {
	return b.latency
}
