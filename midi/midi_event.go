package midi

import (
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

func (m midiEvent) NoteChangesDo(callback func(core.NoteChange)) {
	for _, each := range m.which {
		callback(core.NewNoteChange(m.onoff == noteOn, each, m.velocity))
	}
}

func (m midiEvent) Handle(tim *core.Timeline, when time.Time) {
	// TODO not sure if the noteOn check is correct
	if m.mustHandle != nil && m.onoff == noteOn && !m.mustHandle() {
		return
	}
	status := m.onoff | int64(m.channel-1)
	for _, each := range m.which {
		if err := m.out.WriteShort(status, each, m.velocity); err != nil {
			notify.Errorf("failed to write MIDI data, error:%v", err)
		}
	}
	if m.echoString != "" {
		m.log(status, when)
	}
}

func (m midiEvent) log(status int64, when time.Time) {
	onoff := "on"
	if m.onoff == noteOff {
		onoff = "off"
	}
	fmt.Fprintf(notify.Console.StandardOut, "%s dev=%d ch=%d seq='%s' %s=%d,%v,%d\n",
		when.Format("04:05.000"), m.device, m.channel, m.echoString, onoff, status, m.which, m.velocity)
}

func (m midiEvent) asNoteoff() midiEvent {
	m.onoff = noteOff
	return m
}

type restEvent struct {
	echoString string
	mustHandle core.Condition
}

func (r restEvent) NoteChangesDo(block func(core.NoteChange)) {}

func (r restEvent) Handle(tim *core.Timeline, when time.Time) {}
