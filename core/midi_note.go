package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/emicklei/melrose/notify"
)

type MIDINote struct {
	duration HasValue // faction number or number in milliseconds or time.Duration
	number   HasValue
	velocity HasValue
}

// ToNote() is part of NoteConvertable
func (m MIDINote) ToNote() (Note, error) {
	nr := Int(m.number)
	velocity := Int(m.velocity)
	// check for fraction
	i := getInt(m.duration, true) // do not warn if not integer
	if i == 1 ||
		i == 2 ||
		i == 4 ||
		i == 8 ||
		i == 16 {
		fraction := 1.0 / float32(i)
		return MIDItoNote(fraction, nr, velocity)
	}
	// use as time.Duration or milliseconds
	n, err := MIDItoNote(0.25, nr, velocity)
	if err != nil {
		return Rest4, err
	}
	// 0.25 will be discarded
	n.duration = Duration(m.duration)
	if n.duration < 0 {
		return Rest4, errors.New("MIDI duration cannot be < 0")
	}
	return n, nil
}

func (m MIDINote) S() Sequence {
	n, err := m.ToNote()
	if err != nil {
		notify.Console.Errorf("MIDI to sequence failed:%v", err)
		return EmptySequence
	}
	return n.S()
}

func NewMIDI(duration HasValue, number HasValue, velocity HasValue) MIDINote {
	return MIDINote{duration: duration, number: number, velocity: velocity}
}

func (m MIDINote) Storex() string {
	if s, ok := m.duration.(Storable); ok {
		return fmt.Sprintf("midi(%s,%v,%v)", s.Storex(), m.number, m.velocity)
	}
	return fmt.Sprintf("midi(%v,%v,%v)", m.duration, m.number, m.velocity)
}

func (m MIDINote) Inspect(i Inspection) {
	n, err := m.ToNote()
	if err != nil {
		i.Properties["error"] = err.Error()
		return
	}
	i.Properties["note"] = n.String()
	if d, ok := n.NonFractionBasedDuration(); ok {
		i.Properties["duration"] = d
	} else {
		wholeNoteDuration := WholeNoteDuration(i.Context.Control().BPM())
		i.Properties["duration"] = time.Duration(float32(wholeNoteDuration) * n.DurationFactor())
	}
	i.Properties["velocity"] = m.velocity
}

func IsBlackKey(nr int) bool {
	// https://www.inspiredacoustics.com/en/MIDI_note_numbers_and_center_frequencies
	i := (nr - 21) % 12
	return i == 1 || i == 4 || i == 6 || i == 9 || i == 11
}
