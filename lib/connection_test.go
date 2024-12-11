package lib

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

type readReturn struct {
	n    int
	data []byte
	err  error
}
type writeReturn struct {
	n   int
	err error
}

type mockConn struct {
	readCount  *int
	writeCount *int
	closeCount *int
	readRes    []readReturn
	writeRes   []writeReturn
	closeRes   []error
}

func (m mockConn) Read(p []byte) (int, error) {
	invocation := *m.readCount
	*m.readCount++
	res := m.readRes[invocation%len(m.readRes)]
	copy(p, res.data)
	return res.n, res.err
}

func (m mockConn) Write(p []byte) (int, error) {
	invocation := *m.writeCount
	*m.writeCount++
	res := m.writeRes[invocation%len(m.writeRes)]
	return res.n, res.err
}

func (m mockConn) Close() error {
	invocation := *m.closeCount
	*m.closeCount++
	res := m.closeRes[invocation%len(m.closeRes)]
	return res
}

type mockLatencyGenerator struct {
	delay time.Duration
}

func (m *mockLatencyGenerator) generateLatency(when time.Time) time.Duration {
	return m.delay
}

func TestReadFromSrc(t *testing.T) {
	readCnt := new(int)
	writeCtn := new(int)
	closeCnt := new(int)
	mockSrc := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		readRes: []readReturn{
			{10, []byte("testdata12jibberish"), nil},
			{0, []byte(""), errors.New("some-error")},
		},
	}

	delayQueue := make(chan transitBuffer, 10)
	done := make(chan error, 3)

	c := &connection{
		srcConn:    mockSrc,
		bufferSize: 20,
		latencyGen: &mockLatencyGenerator{time.Millisecond * 2},
		delayQueue: delayQueue,
		done:       done,
		log:        hclog.NewNullLogger(),
	}

	c.readFromSrc()

	transitBuff := <-delayQueue
	err := <-done

	assert.Equal(t, 2, *mockSrc.readCount)
	assert.Equal(t, []byte("testdata12"), transitBuff.data)
	assert.EqualError(t, err, "Error reading data from client some-error")
}

func TestReadFromDest(t *testing.T) {
	readCnt := new(int)
	writeCtn := new(int)
	closeCnt := new(int)
	mockDest := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		readRes: []readReturn{
			{10, []byte("testdata12jibberish"), nil},
			{0, []byte(""), errors.New("some-error")},
		},
	}

	readCnt = new(int)
	writeCtn = new(int)
	closeCnt = new(int)
	mockSrc := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		writeRes: []writeReturn{
			{10, nil},
			{0, errors.New("other-error")},
		},
	}

	done := make(chan error, 3)

	c := &connection{
		srcConn:    mockSrc,
		destConn:   mockDest,
		bufferSize: 20,
		latencyGen: &mockLatencyGenerator{time.Millisecond * 2},
		done:       done,
	}

	c.readFromDest()

	err := <-done

	assert.Equal(t, 2, *mockDest.readCount)
	assert.Equal(t, 1, *mockSrc.writeCount)
	assert.EqualError(t, err, "Error reading data from proxy destination: some-error")
}

func TestReadFromDestSrcWriteError(t *testing.T) {
	readCnt := new(int)
	writeCtn := new(int)
	closeCnt := new(int)
	mockDest := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		readRes: []readReturn{
			{10, []byte("testdata12jibberish"), nil},
			{10, []byte("testdata34jibberish"), nil},
		},
	}

	readCnt = new(int)
	writeCtn = new(int)
	closeCnt = new(int)
	mockSrc := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		writeRes: []writeReturn{
			{10, nil},
			{0, errors.New("other-error")},
		},
	}

	done := make(chan error, 3)

	c := &connection{
		srcConn:    mockSrc,
		destConn:   mockDest,
		bufferSize: 20,
		latencyGen: &mockLatencyGenerator{time.Millisecond * 2},
		done:       done,
	}

	c.readFromDest()

	err := <-done

	assert.Equal(t, 2, *mockDest.readCount)
	assert.Equal(t, 2, *mockSrc.writeCount)
	assert.EqualError(t, err, "Error writing data back to proxy client: other-error")
}

func TestReadFromDelayQueue(t *testing.T) {
	readCnt := new(int)
	writeCtn := new(int)
	closeCnt := new(int)
	mockDest := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		writeRes: []writeReturn{
			{10, nil},
			{0, errors.New("write-error")},
		},
	}

	delayQueue := make(chan transitBuffer, 10)
	done := make(chan error, 3)

	c := &connection{
		destConn:   mockDest,
		bufferSize: 20,
		delayQueue: delayQueue,
		done:       done,
		log:        hclog.NewNullLogger(),
	}

	delayQueue <- transitBuffer{[]byte("testdata"), time.Now().Add(time.Millisecond)}
	delayQueue <- transitBuffer{[]byte("testdata"), time.Now().Add(time.Millisecond * 2)}

	c.readFromDelayQueue()

	err := <-done

	assert.Equal(t, 2, *mockDest.writeCount)
	assert.EqualError(t, err, "Error writing from delay queue to proxy destination: write-error")
}

func TestStart(t *testing.T) {
	readCnt := new(int)
	writeCtn := new(int)
	closeCnt := new(int)
	mockDest := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		readRes: []readReturn{
			{10, []byte("testdata12jibberish"), nil},
			{10, []byte("testdata34jibberish"), nil},
			{10, []byte("testdata56jibberish"), nil},
		},
		writeRes: []writeReturn{
			{10, nil},
			{0, errors.New("dest-write-err")},
		},
		closeRes: []error{nil},
	}

	readCnt = new(int)
	writeCtn = new(int)
	closeCnt = new(int)
	mockSrc := mockConn{
		readCount:  readCnt,
		writeCount: writeCtn,
		closeCount: closeCnt,
		readRes: []readReturn{
			{10, []byte("testdata12jibberish"), nil},
			{10, []byte("testdata34jibberish"), nil},
		},
		writeRes: []writeReturn{
			{10, nil},
		},
		closeRes: []error{nil},
	}

	delayQueue := make(chan transitBuffer, 10)
	done := make(chan error, 3)

	c := &connection{
		srcConn:    mockSrc,
		destConn:   mockDest,
		bufferSize: 20,
		latencyGen: &mockLatencyGenerator{time.Millisecond * 10},
		delayQueue: delayQueue,
		done:       done,
		ctx:        context.TODO(),
		log:        hclog.NewNullLogger(),
	}

	c.start()

	time.Sleep(time.Millisecond * 30)

	assert.Equal(t, 2, *mockDest.writeCount)
	assert.Equal(t, 1, *mockDest.closeCount)
	assert.Equal(t, 1, *mockSrc.closeCount)
}

func TestNewProxyConnectionError(t *testing.T) {
	localAddr, _ := net.ResolveTCPAddr("tcp", ":8000")
	destAddr, _ := net.ResolveTCPAddr("tcp", "nope:3000")

	mockClientConn := mockConn{}

	_, err := newProxyConnection(
		context.TODO(),
		mockClientConn,
		localAddr,
		destAddr,
		0xffff,
		100,
		&mockLatencyGenerator{time.Millisecond * 10},
		hclog.Default(),
		false,
	)

	assert.NotNil(t, err)
}
