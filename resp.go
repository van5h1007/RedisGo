package main

import "bufio"

type Value struct {
	typ   byte
	bulk  string
	array []Value
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')

	if err != nil {
		return "", err
	}

	return line[:len(line)-2], nil
}

func readBulkString(r *bufio.Reader) (Value, error) 

}

func readArray(r *bufio.Reader) (Value, error) {

}