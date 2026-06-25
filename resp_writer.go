package main

import (
	"fmt"
	"io"
)

func writeSimpleString(w io.Writer, s string) error {
	_, err := fmt.Fprintf(w, "+%s\r\n", s)

	return err
}

func writeError(w io.Writer, msg string) error {
	_, err := fmt.Fprintf(w, "-%s\r\n", msg)

	return err
}

func writeBulkString(w io.Writer, s string) error {
	_, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(s), s)

	return err
}

func writeNullBulkString(w io.Writer) error {
	_, err := fmt.Fprintf(w, "$-1\r\n")
	return err
}
