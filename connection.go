package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type transitBuffer struct {
	data       []byte
	delayUntil time.Time
}

type connection struct {
	id                int
	srcConn, destConn io.ReadWriteCloser
	bufferSize        int
	latencyGen        LatencyGenerator
	delayQueue        chan transitBuffer
	done              chan error
	ctx               context.Context
}

func (c *connection) readFromSrc() {
	for {
		buffer := make([]byte, c.bufferSize)
		bytes, err := c.srcConn.Read(buffer)
		receivedAt := time.Now()
		if err != nil {
			c.done <- fmt.Errorf("Error reading data from client %s", err)
			return
		}
		trimmedBuffer := buffer[:bytes]
		desiredLatency := c.latencyGen.generateLatency(receivedAt)
		delayUntil := receivedAt.Add(desiredLatency)

		t := transitBuffer{
			data:       trimmedBuffer,
			delayUntil: delayUntil,
		}

		c.delayQueue <- t

	}
}

func (c *connection) readFromDest() {
	buffer := make([]byte, c.bufferSize)
	for {
		bytes, err := c.destConn.Read(buffer)
		if err != nil {
			c.done <- fmt.Errorf("Error reading data from proxy destination: %s", err)
			return
		}
		trimmedBuffer := buffer[:bytes]

		bytes, err = c.srcConn.Write(trimmedBuffer)
		if err != nil {
			c.done <- fmt.Errorf("Error writing data back to proxy client: %s", err)
			return
		}
	}
}

func (c *connection) readFromDelayQueue() {
	for {
		t := <-c.delayQueue

		time.Sleep(time.Until(t.delayUntil))

		_, err := c.destConn.Write(t.data)
		if err != nil {
			c.done <- fmt.Errorf("Error writing from delay queue to proxy destination: %s", err)
			return
		}
	}
}

// start launches 3 goroutines responsible for handling a proxy connection
// (dest->src, src->queue, queue->dest). This operation will block until
// either an error is sent via the done channel or the context is cancelled.
func (c *connection) start() {
	go c.readFromDest()
	go c.readFromSrc()
	go c.readFromDelayQueue()
	for {
		select {
		case err := <-c.done:
			c.handleError(err)
			return
		case <-c.ctx.Done():
			c.handleStop()
			return
		}
	}
}

func (c *connection) handleError(err error) {
	if !strings.HasSuffix(err.Error(), io.EOF.Error()) {
		fmt.Printf("Closing proxy connection due to an unexpected error: %s\n", err)
	} else {
		fmt.Println("Closing proxy connection (EOF)")
	}
	c.closeProxyConnections()
}

func (c *connection) handleStop() {
	fmt.Printf("Stopping proxy connection #%d\n", c.id)
	c.closeProxyConnections()
}

func (c *connection) closeProxyConnections() {
	c.srcConn.Close()
	c.destConn.Close()
}

func newProxyConnection(
	ctx context.Context,
	clientConn io.ReadWriteCloser,
	srcAddr *net.TCPAddr,
	destAddr *net.TCPAddr,
	bufferSize int,
	latencyGen LatencyGenerator,
	id int,
) (*connection, error) {
	destConn, err := net.DialTCP("tcp", nil, destAddr)
	if err != nil {
		return nil, fmt.Errorf("Error dialing remote address: %s", err)
	}
	c := &connection{
		id:         id,
		srcConn:    clientConn,
		destConn:   destConn,
		bufferSize: bufferSize,
		latencyGen: latencyGen,
		delayQueue: make(chan transitBuffer, 1024),
		done:       make(chan error, 3),
		ctx:        ctx,
	}

	return c, nil
}
