package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadValue(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	r:= bufio.NewReader(strings.NewReader(input))

	v, err:= ReadValue(r)

	if err!= nil {
		t.Fatal(err)
	}

	if v.typ !='*' || len(v.array)!=3 {
		t.Fatalf("expected len of 3 but got %+v", v)
	}

	want := []string{"SET", "foo", "bar"}

	for i,w := range want {
		if v.array[i].bulk != w {
			t.Errorf("element %d: got %q, want %q", i, v.array[i].bulk, w)
		}
	}
}
