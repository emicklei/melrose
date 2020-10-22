package core

import (
	"errors"
	"fmt"
	"time"
)

// noteMidiOffsets maps a tone index (C=0) to the number of semitones on the scale
var noteMidiOffsets = []int{0, 2, 4, 5, 7, 9, 11}

const (
	// maps a tone to an index (C=0)
	nonRestNoteNames = "CDEFGAB"
)

var noteNameToOffset = map[string]int{
	"C": 0,
	"D": 2,
	"E": 4,
	"F": 5,
	"G": 7,
	"A": 9,
	"B": 11,
}

func (n Note) MIDI() int {
	// http://en.wikipedia.org/wiki/Musical_Note
	// C4 = 60 (scientific pitch notation)
	if !n.IsHearable() {
		return 0
	}
	nameOffset := noteNameToOffset[n.Name]
	return ((1 + n.Octave) * 12) + nameOffset + n.Accidental
}

func DurationToFraction(bpm float64, d time.Duration) float32 {
	one := WholeNoteDuration(bpm)
	abs := func(i int) int {
		if i < 0 {
			return -i
		}
		return i
	}
	numbers := []struct {
		fraction float32
		ms       int
	}{
		{1.0, int(one)},
		{0.5, int(one / time.Duration(2))},
		{0.25, int(one / time.Duration(4))},
		{0.125, int(one / time.Duration(8))},
		{0.0625, int(one / time.Duration(16))},
	}
	millis := int(d)
	distance := abs(numbers[0].ms - millis)
	idx := 0
	for c := 1; c < len(numbers); c++ {
		cdistance := abs(numbers[c].ms - millis)
		if cdistance < distance {
			idx = c
			distance = cdistance
		}
	}
	return numbers[idx].fraction
}

func MIDItoNote(fraction float32, nr int, vel int) (Note, error) {
	if fraction < 0 {
		return Rest4, errors.New("MIDI fraction cannot be < 0")
	}
	if nr < 0 || nr > 127 {
		return Rest4, errors.New("MIDI number must be in [0..127]")
	}
	if vel < 0 || vel > 127 {
		return Rest4, errors.New("MIDI velocity must be in [0..127]")
	}
	octave := (nr / 12) - 1
	nrIndex := nr - ((octave + 1) * 12)
	var offsetIndex, offset int
	for o, each := range noteMidiOffsets {
		if each >= nrIndex {
			offsetIndex = o
			offset = each
			break
		}
	}
	accidental := 0
	if nrIndex != offset {
		accidental = -1
	}
	nn, err := NewNote(string(nonRestNoteNames[offsetIndex]), octave, fraction, accidental, false, vel)
	if err != nil {
		return Rest4, err
	}
	return nn, nil
}

type ChannelSelector struct {
	Target Sequenceable
	Number Valueable
}

func NewChannelSelector(target Sequenceable, channel Valueable) ChannelSelector {
	return ChannelSelector{Target: target, Number: channel}
}

func (c ChannelSelector) S() Sequence {
	return c.Target.S()
}

func (c ChannelSelector) Channel() int {
	return Int(c.Number)
}

func (c ChannelSelector) Storex() string {
	if s, ok := c.Target.(Storable); ok {
		return fmt.Sprintf("channel(%v,%s)", c.Number, s.Storex())
	}
	return ""
}
