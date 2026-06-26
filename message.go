package main

import "net"

type Message struct {
	cmd  Command
	conn net.Conn
}
