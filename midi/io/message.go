package io

import (
	"bufio"
	"io"
	"log"
)

const (
	noteoff        byte = 0b10000000
	noteon         byte = 0b10010000
	command_mask   byte = 0b11110000
	iscommand_mask byte = 0b10000000
	channel_mask   byte = 0b00001111
	controlchange  byte = 0xB0 // 10110000 , 176
)

type Message struct {
	command    byte
	channel    byte
	parameter1 byte
	parameter2 byte
}

func (m Message) Status() int64 { return int64(m.command | m.channel) }
func (m Message) Data1() int64  { return int64(m.parameter1) }
func (m Message) Data2() int64  { return int64(m.parameter2) }

var noMessage = Message{}

func WriteMessage(status int64, data1 int64, data2 int64, w io.Writer) error {
	data := []byte{byte(status & 0xFF), byte(data1 & 0xFF), byte(data2 & 0xFF)}
	_, err := w.Write(data)
	return err
}

func ReadMessage(r *bufio.Reader) (Message, error) {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return noMessage, err
		}
		switch b & command_mask {
		case noteon:
			note, err := r.ReadByte()
			if err != nil {
				return noMessage, err
			}
			velocity, err := r.ReadByte()
			if err != nil {
				return noMessage, err
			}
			return Message{command: noteon, channel: b & channel_mask, parameter1: note, parameter2: velocity}, nil
		case noteoff:
			note, err := r.ReadByte()
			if err != nil {
				return noMessage, err
			}
			velocity, err := r.ReadByte()
			if err != nil {
				return noMessage, err
			}
			return Message{command: noteoff, channel: b & channel_mask, parameter1: note, parameter2: velocity}, nil
		case controlchange:
			p1, err := r.ReadByte()
			if err != nil {
				return noMessage, err
			}
			p2, err := r.ReadByte()
			if err != nil {
				return noMessage, err
			}
			return Message{command: controlchange, parameter1: p1, parameter2: p2}, nil
		default:
			log.Printf("unknown command: %b (%d)\n", b, b)
			// consume data bytes
			for seenCommand := false; !seenCommand; {
				data, err := r.ReadByte()
				if err != nil {
					return noMessage, err
				}
				if data&iscommand_mask == iscommand_mask {
					r.UnreadByte()
					break
				}
			}
		}
	}
}
