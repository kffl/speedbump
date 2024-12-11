package lib

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
	Base:          time.Millisecond * 5,
	SineAmplitude: time.Duration(0),
	SinePeriod:    time.Minute,
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
		"localhost",
		8000,
		"localhost:1234",
		0xffff,
		100,
		defaultLatencyCfg,
		"WARN",
		false,
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, err)
	assert.Equal(t, 0xffff, s.bufferSize)
}

func TestNewSpeedbumpInvalidHost(t *testing.T) {
	cfg := SpeedbumpCfg{
		"nope",
		8080,
		"localhost:1234",
		0xffff,
		100,
		defaultLatencyCfg,
		"WARN",
		false,
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, s)
	assert.ErrorContains(t, err, "lookup nope")
}

func TestNewSpeedbumpErrorResolvingLocal(t *testing.T) {
	cfg := SpeedbumpCfg{
		"localhost",
		-1,
		"localhost:1234",
		0xffff,
		100,
		defaultLatencyCfg,
		"WARN",
		false,
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, s)
	assert.True(t, strings.HasPrefix(err.Error(), "Error resolving local"))
}

func TestNewSpeedbumpErrorResolvingDest(t *testing.T) {
	cfg := SpeedbumpCfg{
		"localhost",
		8000,
		"nope:1234",
		0xffff,
		100,
		defaultLatencyCfg,
		"WARN",
		false,
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, s)
	assert.True(t, strings.HasPrefix(err.Error(), "Error resolving destination"))
}

func TestNewSpeedbumpDefaultQueueSize(t *testing.T) {
	cfg := SpeedbumpCfg{
		Port:       8000,
		DestAddr:   "localhost:1234",
		BufferSize: 0xffff,
		// QueueSize is ommitted
		Latency:  defaultLatencyCfg,
		LogLevel: "WARN",
	}
	s, err := NewSpeedbump(&cfg)
	assert.Nil(t, err)
	assert.Equal(t, 1024, s.queueSize)
}

func TestStartListenError(t *testing.T) {
	cfg := SpeedbumpCfg{
		"localhost",
		1, // a privileged port
		"localhost:1234",
		0xffff,
		100,
		defaultLatencyCfg,
		"WARN",
		false,
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
		"localhost",
		8000,
		testSrvAddr,
		0xffff,
		100,
		&LatencyCfg{
			Base:          time.Millisecond * 100,
			SineAmplitude: time.Millisecond * 100,
			SinePeriod:    time.Millisecond * 400,
		},
		"WARN",
		false,
	}
	s, err := NewSpeedbump(&cfg)
	s.Start()

	assert.Nil(t, err)

	tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:8000")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

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

func TestShouldDontDelayWhenDisabled(t *testing.T) {
	port := 9006
	testSrvAddr := fmt.Sprintf("localhost:%d", port)

	go startEchoSrv(port)

	cfg := SpeedbumpCfg{
		"localhost",
		8000,
		testSrvAddr,
		0xffff,
		100,
		&LatencyCfg{
			Base: time.Millisecond * 250,
		},
		"WARN",
		true,
	}
	s, err := NewSpeedbump(&cfg)
	s.Start()

	assert.Nil(t, err)

	tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:8000")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	firstOpStart := time.Now()

	conn.Write([]byte("test-string"))
	res := make([]byte, 1024)
	bytes, _ := conn.Read(res)

	firstOpElapsed := time.Since(firstOpStart)

	trimmedRes := res[:bytes]

	assert.Equal(t, []byte("test-string"), trimmedRes)
	assert.Less(t, firstOpElapsed, 1*time.Millisecond)
}

func TestShouldDelayWhenDisabledAndthenSwitchToEnable(t *testing.T) {
	port := 9006
	testSrvAddr := fmt.Sprintf("localhost:%d", port)

	go startEchoSrv(port)

	cfg := SpeedbumpCfg{
		"localhost",
		8000,
		testSrvAddr,
		0xffff,
		100,
		&LatencyCfg{
			Base: time.Millisecond * 250,
		},
		"WARN",
		true,
	}
	s, err := NewSpeedbump(&cfg)
	s.Start()
	s.Enable()

	assert.Nil(t, err)

	tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:8000")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	firstOpStart := time.Now()

	conn.Write([]byte("test-string"))
	res := make([]byte, 1024)
	bytes, _ := conn.Read(res)

	firstOpElapsed := time.Since(firstOpStart)

	trimmedRes := res[:bytes]

	fmt.Printf("*** %+v\n", firstOpElapsed)
	assert.Equal(t, []byte("test-string"), trimmedRes)
	assert.True(t, isDurationCloseTo(time.Millisecond*250, firstOpElapsed, 20))
}
