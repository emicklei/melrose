package core

import (
	"fmt"
)

type MIDI struct {
	number   Valueable
	velocity Valueable
}

// ToNote() is part of NoteConvertable
func (m MIDI) ToNote() Note {
	nr := Int(m.number)
	velocity := Int(m.velocity)
	return MIDItoNote(nr, velocity)
}

func (m MIDI) S() Sequence {
	return m.ToNote().S()
}

func NewMIDI(number Valueable, velocity Valueable) MIDI {
	return MIDI{number: number, velocity: velocity}
}

func (m MIDI) Storex() string {
	return fmt.Sprintf("midi(%v,%v)", m.number, m.velocity)
}

func (m MIDI) Inspect(i Inspection) {
	n := m.ToNote()
	i.Properties["note"] = n.Storex()
	i.Properties["velocity"] = n.Velocity
}
