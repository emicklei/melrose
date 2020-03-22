package melrose

func (n Note) Repeated(howMany int) Sequence {
	notes := []Note{}
	for i := 0; i < howMany; i++ {
		notes = append(notes, n)
	}
	return BuildSequence(notes)
}

// Pitched creates a new Note with a pitch by a (positive or negative) number of semi tones
func (n Note) Pitched(howManySemitones int) Note {
	simple := MIDItoNote(n.MIDI() + howManySemitones)
	nn, _ := NewNote(simple.Name, simple.Octave, n.duration, simple.Accidental, n.Dotted)
	return nn
}

func (n Note) Octaved(howmuch int) Note {
	nn, _ := NewNote(n.Name, n.Octave+howmuch, n.duration, n.Accidental, n.Dotted)
	return nn
}