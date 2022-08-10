package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSawtoothLatencySummand(t *testing.T) {
	s := sawtoothLatencySummand{
		amplitude: time.Second,
		period:    time.Minute,
	}

	assert.Equal(t, time.Millisecond*0, s.getLatency(time.Duration(0)))
	assert.Equal(t, time.Millisecond*500, s.getLatency(time.Second*15))
	assert.Equal(t, time.Millisecond*-1000, s.getLatency(time.Second*30))
	assert.Equal(t, time.Millisecond*-800, s.getLatency(time.Second*36))
	assert.Equal(t, time.Millisecond*0, s.getLatency(time.Second*60))
}
