package lib

import (
	"math"
	"time"
)

type squareLatencySummand struct {
	amplitude time.Duration
	period    time.Duration
}

func (s squareLatencySummand) getLatency(elapsed time.Duration) time.Duration {
	return time.Duration(
		(2 * (2 * math.Floor((1 / float64 (s.period)) * float64(elapsed)) - math.Floor(2 * (1 / float64 (s.period)) * float64(elapsed))) + 1) * float64(s.amplitude),
	)
}
