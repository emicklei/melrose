package melrose

import "fmt"

type MIDI struct {
	number   Valueable
	velocity Valueable
}

func (m MIDI) S() Sequence {
	nr := Int(m.number)
	velocity := Int(m.velocity)
	return MIDItoNote(nr, velocity).S()
}

func NewMIDI(number Valueable, velocity Valueable) MIDI {
	return MIDI{number: number, velocity: velocity}
}

func (m MIDI) Storex() string {
	return fmt.Sprintf("midi(%v,%v)", m.number, m.velocity)
}
