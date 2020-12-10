package io

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_read(t *testing.T) {
	data := []byte{
		noteon, 65, 127,
		0b11001100, 1, 2, 3, 127,
		noteoff, 23, 3}
	r := bufio.NewReader(bytes.NewReader(data))
	{
		m, err := ReadMessage(r)
		if err != nil {
			t.Error(err)
		}
		t.Log(m)
	}
	{
		m, err := ReadMessage(r)
		if err != nil {
			t.Error(err)
		}
		t.Log(m)
	}
}
