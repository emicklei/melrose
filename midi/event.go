package midi

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose"
	"github.com/rakyll/portmidi"
)

type midiEvent struct {
	echoString string
	which      []int64
	onoff      int64
	channel    int
	velocity   int64
	out        *portmidi.Stream
}

func (m midiEvent) String() string {
	onoff := "on"
	if m.onoff == noteOff {
		onoff = "off"
	}
	return fmt.Sprintf("[%d] %s %s", m.channel, m.echoString, onoff)
}

func (m midiEvent) Handle(tim *melrose.Timeline, when time.Time) {
	if len(m.echoString) > 0 {
		print(m.echoString)
	}
	for _, each := range m.which {
		m.out.WriteShort(m.onoff|int64(m.channel-1), each, m.velocity)
	}
}

func (m midiEvent) asNoteoff() midiEvent {
	m.onoff = noteOff
	// do not echo OFF
	m.echoString = ""
	return m
}
