package melrose

import "strings"

// noteMidiOffsets maps a tone index (C=0) to the number of semitones on the scale
var noteMidiOffsets = []int{0, 2, 4, 5, 7, 9, 11}

func (n Note) MIDI() int {
	// http://en.wikipedia.org/wiki/Musical_Note
	// C4 = 60 (scientific pitch notation)
	if n.IsRest() { // TODO
		return 0
	}
	nameIndex := strings.Index(NonRestNoteNames, n.Name)
	nameOffset := noteMidiOffsets[nameIndex]
	return ((1 + n.Octave) * 12) + nameOffset + n.Accidental
}

func MIDItoNote(nr int) Note {
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
	nn, _ := NewNote(string(NonRestNoteNames[offsetIndex]), octave, 0.25, accidental, false, 1.0)
	return nn
}
