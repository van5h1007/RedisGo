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

func writeInteger (w io.Writer, n int) error {
	_, err := fmt.Fprintf(w, ":%d\r\n", n)

	return err
}

func writeMap(w io.Writer, m map[string]string) error {
	if _, err := fmt.Fprintf(w, "%%%d\r\n", len(m)); err != nil {
		return err
	}
	for k, v := range m {
		if err := writeSimpleString(w, k); err != nil {
			return err
		}
		if err := writeSimpleString(w, v); err != nil {
			return err
		}
	}
	return nil
}
