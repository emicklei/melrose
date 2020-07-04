package core

import (
	"fmt"
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
	if n.IsRest() { // TODO
		return 0
	}
	nameOffset := noteNameToOffset[n.Name]
	return ((1 + n.Octave) * 12) + nameOffset + n.Accidental
}

// TODO handle duration
func MIDItoNote(duration float32, nr int, vel int) Note {
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
	nn, _ := NewNote(string(nonRestNoteNames[offsetIndex]), octave, duration, accidental, false, vel)
	return nn
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
