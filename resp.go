package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

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

func readBulkString(r *bufio.Reader) (Value, error) {
	lenghtLine, err := readLine(r)

	if err != nil {
		return Value{}, err
	}

	length, err := strconv.Atoi(lenghtLine)

	if err != nil {
		return Value{}, fmt.Errorf("bad bulk string %q: %w", lenghtLine, err)
	}

	buf := make([]byte, length+2)

	if _, err := io.ReadFull(r, buf); err != nil {
		return Value{}, err
	}

	return Value{typ: '$', bulk: string(buf[:length])}, nil

}

func readArray(r *bufio.Reader) (Value, error) {
	countLine, err := readLine(r)

	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(countLine)

	if err != nil {
		return Value{}, fmt.Errorf("bad array length %q: %w", countLine, err)
	}

	arr := make([]Value, count)

	for i := 0; i < count; i++ {
		v, err := ReadValue(r)
		if err != nil {
			return Value{}, err
		}

		arr[i] = v
	}
	return Value{typ: '*', array: arr}, nil
}

func ReadValue(r *bufio.Reader) (Value, error) {
	typeByte, err:= r.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeByte{
	case '*':
		return readArray(r)
	case '$':
		return readBulkString(r)
	default:
		return Value{}, fmt.Errorf("unknown resp byte type : %q", typeByte)
	}
}
