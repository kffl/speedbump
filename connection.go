package main

import (
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
	srcConn, destConn io.ReadWriteCloser
	bufferSize        int
	latencyGen        *simpleLatencyGenerator
	delayQueue        chan transitBuffer
	done              chan error
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

func (c *connection) start() {
	go c.readFromDest()
	go c.readFromSrc()
	go c.readFromDelayQueue()
	err := <-c.done
	if !strings.HasSuffix(err.Error(), io.EOF.Error()) {
		fmt.Printf("Closing proxy connection due to an unexpected error: %s\n", err)
	} else {
		fmt.Println("Closing proxy connection (EOF)")
	}
	defer c.srcConn.Close()
	defer c.destConn.Close()
}

func newProxyConnection(
	clientConn io.ReadWriteCloser,
	srcAddr *net.TCPAddr,
	destAddr *net.TCPAddr,
	bufferSize int,
	latencyGen *simpleLatencyGenerator,
) (*connection, error) {
	destConn, err := net.DialTCP("tcp", nil, destAddr)
	if err != nil {
		return nil, fmt.Errorf("Error dialing remote address: %s", err)
	}
	c := &connection{
		srcConn:    clientConn,
		destConn:   destConn,
		bufferSize: bufferSize,
		latencyGen: latencyGen,
		delayQueue: make(chan transitBuffer, 100),
		done:       make(chan error, 3),
	}

	return c, nil
}
