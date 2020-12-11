package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/emicklei/melrose/midi/io"

	"github.com/rakyll/portmidi"
)

type UDPToMIDIListener struct {
	outputStream *portmidi.Stream
	connection   net.PacketConn
}

func newUDPToMIDIListener(port int, deviceID int) (*UDPToMIDIListener, error) {
	if *oDebug {
		log.Println("udp listen on", port)
	}
	c, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	if *oDebug {
		log.Println("portmidi open output device", deviceID)
	}
	s, err := portmidi.NewOutputStream(portmidi.DeviceID(deviceID), 1024, 0) // latency param?
	if err != nil {
		c.Close()
		return nil, err
	}
	return &UDPToMIDIListener{
		outputStream: s,
		connection:   c,
	}, nil
}

// start blocks until error
func (l *UDPToMIDIListener) start() {
	for {
		buffer := make([]byte, 128)
		reader := bufio.NewReader(bytes.NewReader(buffer))
		l.connection.ReadFrom(buffer)
		msg, err := io.ReadMessage(reader)
		if err != nil {
			log.Println("\033[1;33maborted\033[0m reading messages, error:", err)
			return
		}
		if *oDebug {
			log.Println("write portmidi", msg)
		}
		l.outputStream.WriteShort(msg.Status(), msg.Data1(), msg.Data2())
	}
}

func (l *UDPToMIDIListener) close() {
	if l.connection != nil {
		l.connection.Close()
	}
	if l.outputStream != nil {
		l.outputStream.Abort()
		l.outputStream.Close()
	}
}
