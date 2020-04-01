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

var velocityFactors = []float32{F_Pianissimo, F_Piano, F_MezzoPiano, 1.0, F_MezzoForte, F_Forte, F_Fortissimo}

// TODO handle duration
func MIDItoNote(nr int, f float32) Note {
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
	velocityFactor := float32(1.0)
	// quantize the velocity
	minDistance := float32(2.0)
	for _, each := range velocityFactors {
		dist := each - f
		if dist < 0 {
			dist = dist * -1
		}
		if dist < minDistance {
			minDistance = dist
			velocityFactor = each
		}
	}
	nn, _ := NewNote(string(NonRestNoteNames[offsetIndex]), octave, 0.25, accidental, false, float32(velocityFactor))
	return nn
}
