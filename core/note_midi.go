package core

import (
	"errors"
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
	name, octave, accidental := MIDIToNoteParts(nr)
	return MakeNote(name, octave, fraction, accidental, false, vel), nil
}

func MIDIToNoteParts(nr int) (name string, octave int, accidental int) {
	octave = (nr / 12) - 1
	nrIndex := nr - ((octave + 1) * 12)
	var offsetIndex, offset int
	for o, each := range noteMidiOffsets {
		if each >= nrIndex {
			offsetIndex = o
			offset = each
			break
		}
	}
	accidental = 0
	if nrIndex != offset {
		accidental = -1
	}
	return string(nonRestNoteNames[offsetIndex]), octave, accidental
}
