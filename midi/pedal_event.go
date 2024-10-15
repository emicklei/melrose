package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi/transport"
	"github.com/emicklei/melrose/notify"
)

type pedalEvent struct {
	goingDown  bool
	channel    int
	out        transport.MIDIOut
	mustHandle core.Condition
}

func (p pedalEvent) NoteChangesDo(block func(core.NoteChange)) {}

func (p pedalEvent) Handle(tim *core.Timeline, when time.Time) {
	if p.mustHandle != nil && !p.mustHandle() {
		return
	}
	// 0 to 63 = Off, 64 to 127 = On
	var onoff int64 = 0
	if p.goingDown {
		onoff = 127
	}
	// MIDI CC 64	Damper Pedal /Sustain Pedal
	status := controlChange | int64(p.channel-1)
	_ = p.out.WriteShort(status, sustainPedal, onoff)
	if notify.IsDebug() {
		msg := "down"
		if !p.goingDown {
			msg = "up"
		}
		notify.Debugf("midi.pedal channel=%d bytes=[%b(%d),%b(%d),%b(%d)] sustain=%s",
			p.channel, status, status, sustainPedal, sustainPedal, onoff, onoff, msg)
	}
}
