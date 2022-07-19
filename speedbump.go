package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Speedbump struct {
	bufferSize        int
	srcAddr, destAddr net.TCPAddr
	listener          *net.TCPListener
	latencyGen        LatencyGenerator
	nextConnId        int
	// active keeps track of proxy connections that are running
	active sync.WaitGroup
	// ctx is used for notifying proxy connections once Stop() is invoked
	ctx       context.Context
	ctxCancel context.CancelFunc
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
		latencyGen: newSimpleLatencyGenerator(time.Now(), cfg.Latency),
	}
	return s, nil
}

func (s *Speedbump) startAcceptLoop() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed") {
				// the listener was closed, which means that Stop() was called
				return
			} else {
				fmt.Println(fmt.Errorf("Error accepting incoming TCP connection: %s", err))
				continue
			}
		}
		p, err := newProxyConnection(
			s.ctx,
			conn,
			&s.srcAddr,
			&s.destAddr,
			s.bufferSize,
			s.latencyGen,
			s.nextConnId,
		)
		if err != nil {
			fmt.Println(fmt.Errorf("Error creating new proxy connection: %s", err))
			conn.Close()
			continue
		}
		fmt.Printf("Starting a new proxy connection #%d\n", s.nextConnId)
		s.nextConnId++
		s.active.Add(1)
		go s.startProxyConnection(p)
	}
}

func (s *Speedbump) startProxyConnection(p *connection) {
	defer s.active.Done()
	// start will block until a proxy connection is closed
	p.start()
}

func (s *Speedbump) Start() error {
	listener, err := net.ListenTCP("tcp", &s.srcAddr)
	if err != nil {
		return fmt.Errorf("Error starting TCP listener: %s", err)
	}
	s.listener = listener

	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.ctxCancel = cancel

	// startAcceptLoop will block until Stop() is called
	s.startAcceptLoop()
	// wait for active proxy connections to be closed
	s.active.Wait()
	return nil
}

func (s *Speedbump) Stop() {
	// close TCP listener so that startAcceptLoop returns
	s.listener.Close()
	// notify all proxy connections
	s.ctxCancel()
}
