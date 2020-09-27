package core

import (
	"fmt"
)

type MIDI struct {
	fraction Valueable // 0.0625,0.125,0.25,0.5,1,2,4,8,16
	number   Valueable
	velocity Valueable
}

// ToNote() is part of NoteConvertable
func (m MIDI) ToNote() Note {
	f := Float(m.fraction)
	nr := Int(m.number)
	velocity := Int(m.velocity)
	if f > 1.0 {
		f = 1.0 / f
	}
	return MIDItoNote(f, nr, velocity)
}

func (m MIDI) S() Sequence {
	return m.ToNote().S()
}

func NewMIDI(fraction Valueable, number Valueable, velocity Valueable) MIDI {
	return MIDI{fraction: fraction, number: number, velocity: velocity}
}

func (m MIDI) Storex() string {
	return fmt.Sprintf("midi(%v,%v,%v)", m.fraction, m.number, m.velocity)
}

func (m MIDI) Inspect(i Inspection) {
	n := m.ToNote()
	i.Properties["note"] = n.Storex()
	i.Properties["fraction"] = n.fraction
	i.Properties["velocity"] = n.Velocity
}
