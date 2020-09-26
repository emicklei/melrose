package midi

import (
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
	if echoMIDISent && len(m.echoString) > 0 {
		print(m.echoString)
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
		notify.Debugf("ch=%d notes=%s state=%s bytes=[%b(%d),%v,%b(%d)]",
			m.channel, m.echoString, onoff, status, status, m.which, m.velocity, m.velocity)
	}
}

func (m midiEvent) asNoteoff() midiEvent {
	m.onoff = noteOff
	// do not echo OFF
	m.echoString = ""
	return m
}
