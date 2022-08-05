package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleLatencyGeneratorWithSine(t *testing.T) {
	start := time.Now()
	g := newSimpleLatencyGenerator(start, &LatencyCfg{
		Base:          time.Second * 3,
		SineAmplitude: time.Second * 2,
		SinePeriod:    time.Second * 8,
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
		Base:         time.Second * 3,
		SawAmplitude: time.Second * 2,
		SawPeriod:    time.Second * 8,
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

func TestSimpleLatencyGeneratorWithTriangle(t *testing.T) {
	start := time.Now()
	g := newSimpleLatencyGenerator(start, &LatencyCfg{
		Base:              time.Second * 3,
		TriangleAmplitude: time.Second * 2,
		TrianglePeriod:    time.Second * 8,
	})

	startingVal := g.generateLatency(start)
	after2Sec := g.generateLatency(start.Add(time.Second * 2))
	after4Sec := g.generateLatency(start.Add(time.Second * 4))
	after6Sec := g.generateLatency(start.Add(time.Second * 6))
	after2Periods := g.generateLatency(start.Add(time.Second * 16))
	assert.Equal(t, time.Second*3, startingVal)
	assert.Equal(t, time.Second*5, after2Sec)
	assert.Equal(t, time.Second*3, after4Sec)
	assert.Equal(t, time.Second*1, after6Sec)
	assert.Equal(t, time.Second*3, after2Periods)
}
