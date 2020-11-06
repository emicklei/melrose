package midi

import (
	"bytes"
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type midiEvent struct {
	echoString string
	which      []int64
	onoff      int64
	channel    int
	velocity   int64
	out        MIDIOut
}

func (m midiEvent) Handle(tim *core.Timeline, when time.Time) {
	if len(m.echoString) > 0 {
		fmt.Fprintf(notify.Console.DeviceOut, " %s", m.echoString)
	}
	status := m.onoff | int64(m.channel-1)
	for _, each := range m.which {
		m.out.WriteShort(status, each, m.velocity)
	}
	if core.IsDebug() {
		m.log(status, when)
	}
}

func (m midiEvent) log(status int64, when time.Time) {
	onoff := "on"
	if m.onoff == noteOff {
		onoff = "off"
	}
	var echos bytes.Buffer
	for _, each := range m.which {
		n, _ := core.MIDItoNote(0.25, int(each), core.Normal)
		fmt.Fprintf(&echos, "%s ", n.String())
	}
	fmt.Fprintf(notify.Console.StandardOut, "midi.note: t=%s ch=%d %s %s d=[%d,%v,%d]\n",
		when.Format("04:05.000"), m.channel, echos.String(), onoff, status, m.which, m.velocity)
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
	if len(r.echoString) > 0 {
		fmt.Fprintf(notify.Console.DeviceOut, " %s", r.echoString)
	}
}
