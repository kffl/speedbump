package lib

import "time"

type sawtoothLatencySummand struct {
	amplitude time.Duration
	period    time.Duration
}

func (s sawtoothLatencySummand) getLatency(elapsed time.Duration) time.Duration {
	return time.Duration(
		(float64((elapsed+s.period/2)%s.period)/float64(s.period))*float64(s.amplitude*2),
	) - s.amplitude
}
