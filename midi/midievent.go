package midi

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"

	"github.com/rakyll/portmidi"
)

var echoMIDISent bool

type midiEvent struct {
	echoString string
	which      []int64
	onoff      int64
	channel    int
	velocity   int64
	out        *portmidi.Stream
}

func (m midiEvent) Handle(tim *core.Timeline, when time.Time) {
	// if echoMIDISent && len(m.echoString) > 0 {
	if len(m.echoString) > 0 {
		fmt.Fprintf(notify.Console.DeviceOut, m.echoString)
	}
	status := m.onoff | int64(m.channel-1)
	for _, each := range m.which {
		m.out.WriteShort(status, each, m.velocity)
	}
	if core.IsDebug() {
		onoff := "on"
		if m.onoff == noteOff {
			onoff = "off"
		}
		fmt.Fprintf(notify.Console.StandardOut, "ch=%d notes=%s state=%s bytes=[%b(%d),%v,%b(%d)]\n",
			m.channel, m.echoString, onoff, status, status, m.which, m.velocity, m.velocity)
	}
}

func (m midiEvent) asNoteoff() midiEvent {
	m.onoff = noteOff
	// do not echo OFF
	m.echoString = ""
	return m
}

type restEvent struct {
	echoString string
}

func (r restEvent) Handle(tim *core.Timeline, when time.Time) {
	if echoMIDISent && len(r.echoString) > 0 {
		fmt.Fprintf(notify.Console.DeviceOut, r.echoString)
	}
}
