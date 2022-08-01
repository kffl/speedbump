package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSquareLatencySummand(t *testing.T) {
	s := squareLatencySummand{
		amplitude: time.Second * 2,
		period:    time.Minute,
	}

	assert.Equal(t, s.getLatency(time.Duration(0)), time.Second * 2)
	assert.Equal(t, s.getLatency(time.Second * 15), time.Second * 2)
	assert.Equal(t, s.getLatency(time.Second * 30), time.Second * -2)
	assert.Equal(t, s.getLatency(time.Second * 45), time.Second * -2)
	assert.Equal(t, s.getLatency(time.Second * 60), time.Second * 2)
	assert.Equal(t, s.getLatency(time.Second * 84), time.Second * 2)
	assert.Equal(t, s.getLatency(time.Second * 90), time.Second * -2)
}
