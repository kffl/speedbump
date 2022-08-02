package lib

import "time"


type squareLatencySummand struct {
	amplitude time.Duration
	period    time.Duration
}

func (s squareLatencySummand) getLatency(elapsed time.Duration) time.Duration {
	return time.Duration(
		(4 * (elapsed / s.period) - 2 * ((2 * elapsed) / s.period) + 1) * s.amplitude,
	)
}
