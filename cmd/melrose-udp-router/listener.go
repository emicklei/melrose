package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/emicklei/melrose/midi/io"

	"github.com/rakyll/portmidi"
)

type UDPToMIDIListener struct {
	outputStream *portmidi.Stream
	connection   net.Conn
	listener     net.Listener
}

func newUDPToMIDIListener(host string, port int, deviceID int) (*UDPToMIDIListener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	c, err := lis.Accept()
	if err != nil {
		return nil, err
	}
	s, err := portmidi.NewOutputStream(portmidi.DeviceID(deviceID), 1024, 0) // latency param?
	if err != nil {
		c.Close()
		return nil, err
	}
	return &UDPToMIDIListener{
		outputStream: s,
		connection:   c,
		listener:     lis,
	}, nil
}

// start blocks until error
func (l *UDPToMIDIListener) start() {
	reader := bufio.NewReader(l.connection)
	for {
		msg, err := io.ReadMessage(reader)
		if err != nil {
			log.Println("aborted reading messages, error:", err)
			return
		}
		l.outputStream.WriteShort(msg.Status(), msg.Data1(), msg.Data2())
	}
}

func (l *UDPToMIDIListener) close() {
	if l.connection != nil {
		l.connection.Close()
	}
	if l.listener != nil {
		l.listener.Close()
	}
	if l.outputStream != nil {
		l.outputStream.Abort()
		l.outputStream.Close()
	}
}
