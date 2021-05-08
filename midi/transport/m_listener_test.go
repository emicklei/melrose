package transport

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

type noteCollector struct {
	noteOn        bool
	noteOff       bool
	controlChange bool
	channel       int
	number        int
	data2         int
}

func (n *noteCollector) ControlChange(channel, number, value int) {
	n.channel = channel
	n.number = number
	n.data2 = value
	n.controlChange = true
}

func (n *noteCollector) NoteOn(channel int, note core.Note) {
	n.channel = channel
	n.number = note.MIDI()
	n.data2 = note.Velocity
	n.noteOn = true
}
func (n *noteCollector) NoteOff(channel int, note core.Note) {
	n.channel = channel
	n.number = note.MIDI()
	n.data2 = note.Velocity
	n.noteOff = true
}

func Test_mListener_HandleMIDIMessage(t *testing.T) {
	nc := new(noteCollector)
	lis := newMListener()
	lis.Add(nc)
	ch := 4
	status := noteOn | int16(ch-1)
	nr := 85
	v := 50
	lis.HandleMIDIMessage(status, nr, v)
	if got, want := nc.channel, ch; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := nc.noteOn, true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := nc.data2, 50; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
