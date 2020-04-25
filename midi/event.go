package midi

import (
	"time"

	"github.com/emicklei/melrose"
	"github.com/rakyll/portmidi"
)

type midiEvent struct {
	which    []int64
	onoff    int
	channel  int
	velocity int64
	out      *portmidi.Stream
}

func (m midiEvent) Handle(tim *melrose.Timeline, when time.Time) {
	for _, each := range m.which {
		m.out.WriteShort(int64(m.onoff|(m.channel-1)), each, m.velocity)
	}
}

func (m midiEvent) asNoteoff() midiEvent {
	m.onoff = noteOff
	return m
}
