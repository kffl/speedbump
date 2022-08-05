package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTriangleLatencySummand(t *testing.T) {
	s := triangleLatencySummand{
		amplitude: time.Second,
		period:    time.Minute,
	}

	assert.Equal(t, time.Millisecond*0, s.getLatency(time.Duration(0)))
	assert.Equal(t, time.Millisecond*500, s.getLatency(time.Millisecond*7500))
	assert.Equal(t, time.Millisecond*1000, s.getLatency(time.Millisecond*15000))
	assert.Equal(t, time.Millisecond*500, s.getLatency(time.Millisecond*22500))
	assert.Equal(t, time.Millisecond*0, s.getLatency(time.Millisecond*30000))
	assert.Equal(t, -time.Millisecond*500, s.getLatency(time.Millisecond*37500))
	assert.Equal(t, -time.Millisecond*1000, s.getLatency(time.Millisecond*45000))
	assert.Equal(t, -time.Millisecond*500, s.getLatency(time.Millisecond*52500))
	assert.Equal(t, time.Millisecond*0, s.getLatency(time.Millisecond*60000))
}
