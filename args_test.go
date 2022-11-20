package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseArgsDefault(t *testing.T) {
	cfg, err := parseArgs([]string{"localhost:80"})
	assert.Nil(t, err)
	assert.Equal(t, cfg.DestAddr, "localhost:80")
	assert.Equal(t, cfg.Port, 8000)
	assert.Equal(t, 0xffff+1, cfg.BufferSize)
	assert.Equal(t, time.Millisecond*5, cfg.Latency.Base)
	assert.Equal(t, time.Duration(0), cfg.Latency.SineAmplitude)
}

func TestParseArgsError(t *testing.T) {
	_, err := parseArgs([]string{"--nope", "localhost:80"})
	assert.NotNil(t, err)
}

func TestParseArgsAll(t *testing.T) {
	cfg, err := parseArgs(
		[]string{
			"--host=somehost",
			"--port=1234",
			"--buffer=200B",
			"--latency=100ms",
			"--sine-amplitude=50ms",
			"--sine-period=1m",
			"--square-amplitude=123ms",
			"--square-period=3m",
			"--triangle-amplitude=150ms",
			"--triangle-period=2m",
			"host:777",
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, cfg.DestAddr, "host:777")
	assert.Equal(t, cfg.Host, "somehost")
	assert.Equal(t, cfg.Port, 1234)
	assert.Equal(t, 200, cfg.BufferSize)
	assert.Equal(t, time.Millisecond*100, cfg.Latency.Base)
	assert.Equal(t, time.Millisecond*50, cfg.Latency.SineAmplitude)
	assert.Equal(t, time.Minute, cfg.Latency.SinePeriod)
	assert.Equal(t, time.Duration(0), cfg.Latency.SawAmplitude)
	assert.Equal(t, time.Duration(0), cfg.Latency.SawPeriod)
	assert.Equal(t, time.Millisecond*123, cfg.Latency.SquareAmplitude)
	assert.Equal(t, time.Minute*3, cfg.Latency.SquarePeriod)
	assert.Equal(t, time.Millisecond*150, cfg.Latency.TriangleAmplitude)
	assert.Equal(t, time.Minute*2, cfg.Latency.TrianglePeriod)
}
