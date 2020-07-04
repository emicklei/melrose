package core

import (
	"fmt"
)

type MIDI struct {
	duration Valueable // 0.0625,0.125,0.25,0.5,1,2,4,8,16
	number   Valueable
	velocity Valueable
}

// ToNote() is part of NoteConvertable
func (m MIDI) ToNote() Note {
	dur := Float(m.duration)
	nr := Int(m.number)
	velocity := Int(m.velocity)
	if dur > 1.0 {
		dur = 1.0 / dur
	}
	return MIDItoNote(dur, nr, velocity)
}

func (m MIDI) S() Sequence {
	return m.ToNote().S()
}

func NewMIDI(duration Valueable, number Valueable, velocity Valueable) MIDI {
	return MIDI{duration: duration, number: number, velocity: velocity}
}

func (m MIDI) Storex() string {
	return fmt.Sprintf("midi(%v,%v,%v)", m.duration, m.number, m.velocity)
}

func (m MIDI) Inspect(i Inspection) {
	n := m.ToNote()
	i.Properties["note"] = n.Storex()
	i.Properties["duration"] = n.duration
	i.Properties["velocity"] = n.Velocity
}
