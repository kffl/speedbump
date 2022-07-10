package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleLatencyGeneratorWithSine(t *testing.T) {
	start := time.Now()
	g := newSimpleLatencyGenerator(start, &LatencyCfg{
		base:          time.Second * 3,
		sineAmplitude: time.Second * 2,
		sinePeriod:    time.Second * 8,
	})

	startingVal := g.generateLatency(start)
	after2Sec := g.generateLatency(start.Add(time.Second * 2))
	after4Sec := g.generateLatency(start.Add(time.Second * 4))
	after2Periods := g.generateLatency(start.Add(time.Second * 16))
	assert.Equal(t, time.Second*3, startingVal)
	assert.Equal(t, time.Second*5, after2Sec)
	assert.Equal(t, time.Second*3, after4Sec)
	assert.Equal(t, time.Second*3, after2Periods)
}

func TestSimpleLatencyGeneratorWithSawtooth(t *testing.T) {
	start := time.Now()
	g := newSimpleLatencyGenerator(start, &LatencyCfg{
		base:         time.Second * 3,
		sawAmplitute: time.Second * 2,
		sawPeriod:    time.Second * 8,
	})

	startingVal := g.generateLatency(start)
	after2Sec := g.generateLatency(start.Add(time.Second * 2))
	after4Sec := g.generateLatency(start.Add(time.Second * 4))
	after2Periods := g.generateLatency(start.Add(time.Second * 16))
	assert.Equal(t, time.Second*3, startingVal)
	assert.Equal(t, time.Second*4, after2Sec)
	assert.Equal(t, time.Second*1, after4Sec)
	assert.Equal(t, time.Second*3, after2Periods)
}
