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

func (m midiEvent) String() string {
	onoff := "on"
	if m.onoff == noteOff {
		onoff = "off"
	}
	return fmt.Sprintf("ch=%d nrs=%v notes=%s state=%s", m.channel, m.which, m.echoString, onoff)
}

func (m midiEvent) Handle(tim *core.Timeline, when time.Time) {
	if echoMIDISent && len(m.echoString) > 0 {
		print(m.echoString)
	}
	if core.IsDebug() {
		notify.Debugf("%s", m.String())
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