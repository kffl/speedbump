package main

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

	assert.Equal(t, s.getLatency(time.Duration(0)), time.Millisecond*0)
	assert.Equal(t, s.getLatency(time.Second*15), time.Millisecond*500)
	assert.Equal(t, s.getLatency(time.Second*30), time.Millisecond*-1000)
	assert.Equal(t, s.getLatency(time.Second*36), time.Millisecond*-800)
	assert.Equal(t, s.getLatency(time.Second*60), time.Millisecond*0)
}
