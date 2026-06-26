package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	ln    net.Listener
	kv    *KV
	peers map[*Peer]bool

	msgCh     chan Message
	addPeerCh chan *Peer
	delPeerCh chan *Peer
}

func NewServer() *Server {
	return &Server{
		kv:        NewKV(),
		peers:     make(map[*Peer]bool),
		msgCh:     make(chan Message),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
	}
}

func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	log.Println("listening on", addr)
	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				log.Println("handle message error:", err)
			}

		case peer := <-s.addPeerCh:
			log.Println("peer connected:", peer.conn.RemoteAddr())
			s.peers[peer] = true

		case peer := <-s.delPeerCh:
			log.Println("peer disconnected:", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	peer := NewPeer(conn)
	s.addPeerCh <- peer
	defer func() { s.delPeerCh <- peer }()

	r := bufio.NewReader(conn)

	for {
		v, err := ReadValue(r)
		if err == io.EOF {
			log.Println("client disconnected:", conn.RemoteAddr())
			return
		}
		if err != nil {
			log.Println("parse error:", err)
			return
		}

		cmd, err := parseCommand(v)
		if err != nil {
			log.Println("command error:", err)
			if writeErr := writeError(conn, "ERR "+err.Error()); writeErr != nil {
				log.Println("write error:", writeErr)
				return
			}
			continue
		}

		s.msgCh <- Message{cmd: cmd, peer: peer}
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch c := msg.cmd.(type) {

	case PingCommand:
		return writeSimpleString(msg.peer.conn, "PONG")

	case SetCommand:
		s.kv.Set(c.key, c.val)
		return writeSimpleString(msg.peer.conn, "OK")

	case GetCommand:
		val, ok := s.kv.Get(c.key)

		if !ok {
			return writeNullBulkString(msg.peer.conn)
		}
		return writeBulkString(msg.peer.conn, val)

	default:
		return fmt.Errorf("unhandled command type: %T", c)
	}
}
