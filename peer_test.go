package main

import (
	"bufio"
	"net"
	"testing"
)

func TestHelloCommand(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:5001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	conn.Write([]byte("*2\r\n$5\r\nHELLO\r\n$1\r\n2\r\n"))

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
	
	if line[0] != '%' && line[0] != '+' && line[0] != ':' {
		t.Errorf("Unexpected RESP reply: %s", line)
	}
}
