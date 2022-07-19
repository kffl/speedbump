package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var defaultLatencyCfg = &LatencyCfg{
	base:          time.Millisecond * 5,
	sineAmplitude: time.Duration(0),
	sinePeriod:    time.Minute,
}

func startEchoSrv(port int) error {
	srv, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	defer srv.Close()
	for {
		conn, err := srv.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			io.Copy(c, c)
		}(conn)
	}
}

func TestNewSpeedbump(t *testing.T) {
	cfg := SpeedbumpCfg{
		8000,
		"localhost:1234",
		0xffff,
		defaultLatencyCfg,
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, err)
	assert.Equal(t, 0xffff, s.bufferSize)
}

func TestNewSpeedbumpErrorResolvingLocal(t *testing.T) {
	cfg := SpeedbumpCfg{
		-1,
		"localhost:1234",
		0xffff,
		defaultLatencyCfg,
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, s)
	assert.True(t, strings.HasPrefix(err.Error(), "Error resolving local"))
}

func TestNewSpeedbumpErrorResolvingDest(t *testing.T) {
	cfg := SpeedbumpCfg{
		8000,
		"nope:1234",
		0xffff,
		defaultLatencyCfg,
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, s)
	assert.True(t, strings.HasPrefix(err.Error(), "Error resolving destination"))
}

func TestStartListenError(t *testing.T) {
	cfg := SpeedbumpCfg{
		1, // a privileged port
		"localhost:1234",
		0xffff,
		defaultLatencyCfg,
	}
	s, _ := NewSpeedbump(&cfg)

	err := s.Start()

	assert.True(t, strings.HasPrefix(err.Error(), "Error starting TCP listener"))
}

func isDurationCloseTo(expected time.Duration, obtianed time.Duration, percentage int) bool {
	absoluteError := int(expected) - int(obtianed)
	if absoluteError < 0 {
		absoluteError *= -1
	}
	errorPercentage := float64(absoluteError) / float64(expected) * 100.0
	return errorPercentage < float64(percentage)
}

func TestSpeedbumpWithEchoServer(t *testing.T) {
	port := 9006
	testSrvAddr := fmt.Sprintf("localhost:%d", port)

	go startEchoSrv(port)

	cfg := SpeedbumpCfg{
		8000,
		testSrvAddr,
		0xffff,
		&LatencyCfg{
			base:          time.Millisecond * 100,
			sineAmplitude: time.Millisecond * 100,
			sinePeriod:    time.Millisecond * 400,
		},
	}
	s, err := NewSpeedbump(&cfg)
	go s.Start()

	assert.Nil(t, err)

	tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:8000")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)

	firstOpStart := time.Now()

	conn.Write([]byte("test-string"))
	res := make([]byte, 1024)
	bytes, _ := conn.Read(res)

	firstOpElapsed := time.Since(firstOpStart)

	trimmedRes := res[:bytes]

	assert.Equal(t, []byte("test-string"), trimmedRes)
	assert.True(t, isDurationCloseTo(time.Millisecond*100, firstOpElapsed, 20))

	// after ~100ms since test start the added delay will be at 200ms (100ms base + 100ms sine wave max)
	secondOpStart := time.Now()

	conn.Write([]byte("another-test"))
	res = make([]byte, 1024)
	bytes, _ = conn.Read(res)

	secondOpElapsed := time.Since(secondOpStart)

	trimmedRes = res[:bytes]

	s.Stop()

	assert.Equal(t, []byte("another-test"), trimmedRes)
	assert.True(t, isDurationCloseTo(time.Millisecond*200, secondOpElapsed, 20))
}
