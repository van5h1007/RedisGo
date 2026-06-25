package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port:5001")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	for {
		v, err := ReadValue(r)
		if err == io.EOF {
			log.Println("clint disconnected:", conn.RemoteAddr())
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

		if err := handleCommand(conn, cmd); err != nil {
			log.Println("handle error:", err)
			return
		}
	}
}

func handleCommand(conn net.Conn, cmd Command) error {
	switch c := cmd.(type) {
	case PingCommand:
		return writeSimpleString(conn, "PONG")

	case SetCommand:
		fmt.Printf("SET %s = %s\n", c.key, c.val)
		return writeSimpleString(conn, "OK")

	case GetCommand:
		fmt.Printf("GET %s (no store yet, returning nil)\n", c.key)
		return writeSimpleString(conn, "OK")

	default:
		return fmt.Errorf("unhandled command type: %T", cmd)
	}
}
