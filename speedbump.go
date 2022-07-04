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

type SpeedbumpCfg struct {
	Port       int
	DestAddr   string
	BufferSize int
	Latency    *LatencyCfg
}

func NewSpeedbump(cfg *SpeedbumpCfg) (*Speedbump, error) {
	localTCPAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("Error resolving local address: %s", err)
	}
	destTCPAddr, err := net.ResolveTCPAddr("tcp", cfg.DestAddr)
	if err != nil {
		return nil, fmt.Errorf("Error resolving destination address: %s", err)
	}
	s := &Speedbump{
		bufferSize: int(cfg.BufferSize),
		srcAddr:    *localTCPAddr,
		destAddr:   *destTCPAddr,
		latencyGen: &simpleLatencyGenerator{
			start: time.Now(),
			cfg:   cfg.Latency,
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
		go p.start()
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
