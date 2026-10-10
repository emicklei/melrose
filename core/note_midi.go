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

// DurationToFraction returns the note length (and whether it is dotted) that is closest to d at the given tempo.
func DurationToFraction(bpm float64, d time.Duration) (fraction float32, dotted bool) {
	return FractionToDurationParts(float64(d) / float64(WholeNoteDuration(bpm)))
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
