package melrose

// Pitched creates a new Note with a pitch by a (positive or negative) number of semi tones
func (n Note) Pitched(howManySemitones int) Note {
	if howManySemitones == 0 {
		return n
	}
	simple := MIDItoNote(n.MIDI()+howManySemitones, 1.0)
	nn, _ := NewNote(simple.Name, simple.Octave, n.duration, simple.Accidental, n.Dotted, n.velocityFactor)
	return nn
}

func (n Note) Octaved(howmuch int) Note {
	if howmuch == 0 {
		return n
	}
	nn, _ := NewNote(n.Name, n.Octave+howmuch, n.duration, n.Accidental, n.Dotted, n.velocityFactor)
	return nn
}
