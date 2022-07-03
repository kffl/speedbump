package main

import (
	"fmt"
	"net"
	"time"
)

type Speedbump struct {
	bufferSize        int
	srcAddr, destAddr net.TCPAddr
	listener          *net.TCPListener
	latencyGen        *simpleLatencyGenerator
}

func NewSpeedbump(
	port int,
	destAddr string,
	bufferSize int,
	latencyCfg *LatencyCfg,
) (*Speedbump, error) {
	localTCPAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error resolving local address: %s", err)
	}
	destTCPAddr, err := net.ResolveTCPAddr("tcp", destAddr)
	if err != nil {
		return nil, fmt.Errorf("Error resolving destination address: %s", err)
	}
	s := &Speedbump{
		bufferSize: bufferSize,
		srcAddr:    *localTCPAddr,
		destAddr:   *destTCPAddr,
		latencyGen: &simpleLatencyGenerator{
			start: time.Now(),
			cfg:   latencyCfg,
		},
	}
	return s, nil
}

func (s *Speedbump) startAcceptLoop() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			fmt.Println(fmt.Errorf("Error accepting incoming TCP connection: %s", err))
			continue
		}
		p, err := newProxyConnection(conn, &s.srcAddr, &s.destAddr, s.bufferSize, s.latencyGen)
		if err != nil {
			fmt.Println(fmt.Errorf("Error creating new proxy connection: %s", err))
			conn.Close()
			continue
		}
		fmt.Println("Starting a new proxy connection...")
		p.start()
	}
}

func (s *Speedbump) Start() error {
	listener, err := net.ListenTCP("tcp", &s.srcAddr)
	if err != nil {
		return fmt.Errorf("Error starting TCP listener: %s", err)
	}
	s.listener = listener
	s.startAcceptLoop()
	return nil
}
