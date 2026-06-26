package main

import (
	"log"
)

func main() {
	// ln, err := net.Listen("tcp", ":5001")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Listening on port:5001")

	// kv := NewKV()

	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		log.Println("accept error:", err)
	// 		continue
	// 	}
	// 	go handleConn(conn, kv)
	// }
	server := NewServer()

	log.Fatal(server.Start(":5001"))
}

// func handleConn(conn net.Conn, kv *KV) {
// 	defer conn.Close()

// 	r := bufio.NewReader(conn)

// 	for {
// 		v, err := ReadValue(r)
// 		if err == io.EOF {
// 			log.Println("clint disconnected:", conn.RemoteAddr())
// 			return
// 		}
// 		if err != nil {
// 			log.Println("parse error:", err)
// 			return
// 		}

// 		cmd, err := parseCommand(v)

// 		if err != nil {
// 			log.Println("command error:", err)

// 			if writeErr := writeError(conn, "ERR "+err.Error()); writeErr != nil {
// 				log.Println("write error:", writeErr)
// 				return
// 			}
// 			continue
// 		}

// 		if err := handleCommand(conn, cmd, kv); err != nil {
// 			log.Println("handle error:", err)
// 			return
// 		}
// 	}
// }

// func handleCommand(conn net.Conn, cmd Command, kv *KV) error {
// 	switch c := cmd.(type) {
// 	case PingCommand:
// 		return writeSimpleString(conn, "PONG")

// 	case SetCommand:
// 		kv.Set(c.key, c.val)
// 		return writeSimpleString(conn, "OK")

// 	case GetCommand:
// 		val, ok := kv.Get(c.key)

// 		if !ok {
// 			return writeNullBulkString(conn)
// 		}
// 		return writeBulkString(conn, val)

// 	default:
// 		return fmt.Errorf("unhandled command type: %T", cmd)
// 	}
// }
