package lib

import (
	"math"
	"time"
)

type sineLatencySummand struct {
	amplitude time.Duration
	period    time.Duration
}

func (s sineLatencySummand) getLatency(elapsed time.Duration) time.Duration {
	return time.Duration(
		math.Sin(
			float64(elapsed)/float64(s.period)*math.Pi*2,
		) * float64(
			s.amplitude,
		))
}
