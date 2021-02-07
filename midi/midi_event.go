package midi

import (
	"bytes"
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi/transport"
	"github.com/emicklei/melrose/notify"
)

type midiEvent struct {
	echoString string
	which      []int64
	onoff      int64
	channel    int
	velocity   int64
	device     int
	out        transport.MIDIOut
	mustHandle core.Condition
}

func (m midiEvent) Handle(tim *core.Timeline, when time.Time) {
	// TODO not sure if the noteOn check is correct
	if m.mustHandle != nil && m.onoff == noteOn && !m.mustHandle() {
		return
	}
	if len(m.echoString) > 0 {
		fmt.Fprintf(notify.Console.DeviceOut, " %s", m.echoString)
	}
	status := m.onoff | int64(m.channel-1)
	for _, each := range m.which {
		if err := m.out.WriteShort(status, each, m.velocity); err != nil {
			notify.Errorf("failed to write MIDI data, error:%v", err)
		}
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
	for i, each := range m.which {
		if i > 0 {
			fmt.Fprintf(&echos, " ")
		}
		n, _ := core.MIDItoNote(0.25, int(each), core.Normal) // TODO
		fmt.Fprintf(&echos, "%s", n.String())
	}
	fmt.Fprintf(notify.Console.StandardOut, "midi.note: t=%s dev=%d ch=%d seq='%s' %s=%d,%v,%d\n",
		when.Format("04:05.000"), m.device, m.channel, echos.String(), onoff, status, m.which, m.velocity)
}

func (m midiEvent) asNoteoff() midiEvent {
	m.onoff = noteOff
	// do not echo OFF
	m.echoString = ""
	return m
}

type restEvent struct {
	echoString string
	mustHandle core.Condition
}

func (r restEvent) Handle(tim *core.Timeline, when time.Time) {
	if r.mustHandle != nil && !r.mustHandle() {
		return
	}
	if len(r.echoString) > 0 {
		fmt.Fprintf(notify.Console.DeviceOut, " %s", r.echoString)
	}
}
