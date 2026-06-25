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
