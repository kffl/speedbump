package lib

import (
	"math"
	"time"
)

type triangleLatencySummand struct {
	amplitude time.Duration
	period    time.Duration
}

func (t triangleLatencySummand) getLatency(elapsed time.Duration) time.Duration {
	a, p, x := float64(t.amplitude), float64(t.period), float64(elapsed)
	return time.Duration(4*a/p*math.Abs(math.Mod(((math.Mod((x-p/4), p))+p), p)-p/2) - a)
}
