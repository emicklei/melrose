package io

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
)

func TestReadMessage_NoteOn(t *testing.T) {
	data := []byte{noteon, 65, 127}
	r := bufio.NewReader(bytes.NewReader(data))
	m, err := ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	if m.command != noteon {
		t.Errorf("got %v want %v", m.command, noteon)
	}
	if m.channel != 0 {
		t.Errorf("got %v want %v", m.channel, 0)
	}
	if m.parameter1 != 65 {
		t.Errorf("got %v want %v", m.parameter1, 65)
	}
	if m.parameter2 != 127 {
		t.Errorf("got %v want %v", m.parameter2, 127)
	}
}

func TestReadMessage_NoteOff(t *testing.T) {
	data := []byte{noteoff, 23, 3}
	r := bufio.NewReader(bytes.NewReader(data))
	m, err := ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	if m.command != noteoff {
		t.Errorf("got %v want %v", m.command, noteoff)
	}
	if m.channel != 0 {
		t.Errorf("got %v want %v", m.channel, 0)
	}
	if m.parameter1 != 23 {
		t.Errorf("got %v want %v", m.parameter1, 23)
	}
	if m.parameter2 != 3 {
		t.Errorf("got %v want %v", m.parameter2, 3)
	}
}

func TestReadMessage_ControlChange(t *testing.T) {
	data := []byte{controlchange, 2, 3}
	r := bufio.NewReader(bytes.NewReader(data))
	m, err := ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	if m.command != controlchange {
		t.Errorf("got %v want %v", m.command, controlchange)
	}
	if m.parameter1 != 2 {
		t.Errorf("got %v want %v", m.parameter1, 2)
	}
	if m.parameter2 != 3 {
		t.Errorf("got %v want %v", m.parameter2, 3)
	}
}

func TestReadMessage_UnknownCommand(t *testing.T) {
	data := []byte{0b11001100, 1, 2, 3, 127, noteon, 1, 1}
	r := bufio.NewReader(bytes.NewReader(data))
	m, err := ReadMessage(r)
	if err != nil {
		t.Fatal(err)
	}
	if m.command != noteon {
		t.Errorf("got %v want %v", m.command, noteon)
	}
}

func TestReadMessage_ReadByteError(t *testing.T) {
	r := bufio.NewReader(&errReader{errors.New("read error")})
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadMessage_NoteOnReadByteError1(t *testing.T) {
	data := []byte{noteon}
	r := bufio.NewReader(bytes.NewReader(data))
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadMessage_NoteOnReadByteError2(t *testing.T) {
	data := []byte{noteon, 1}
	r := bufio.NewReader(bytes.NewReader(data))
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadMessage_NoteOffReadByteError1(t *testing.T) {
	data := []byte{noteoff}
	r := bufio.NewReader(bytes.NewReader(data))
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadMessage_NoteOffReadByteError2(t *testing.T) {
	data := []byte{noteoff, 1}
	r := bufio.NewReader(bytes.NewReader(data))
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadMessage_ControlChangeReadByteError1(t *testing.T) {
	data := []byte{controlchange}
	r := bufio.NewReader(bytes.NewReader(data))
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadMessage_ControlChangeReadByteError2(t *testing.T) {
	data := []byte{controlchange, 1}
	r := bufio.NewReader(bytes.NewReader(data))
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestReadMessage_UnknownCommandReadByteError(t *testing.T) {
	data := []byte{0b11001100, 1}
	r := bufio.NewReader(bytes.NewReader(data))
	_, err := ReadMessage(r)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestWriteMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	err := WriteMessage(1, 2, 3, buf)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), "\x01\x02\x03"; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
}

func TestWriteMessage_Error(t *testing.T) {
	w := &errWriter{errors.New("write error")}
	err := WriteMessage(1, 2, 3, w)
	if err == nil {
		t.Fatal("expected an error")
	}
}

// errReader is a helper struct for testing io error cases.
type errReader struct {
	err error
}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

// errWriter is a helper struct for testing io error cases.
type errWriter struct {
	err error
}

func (w *errWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}
