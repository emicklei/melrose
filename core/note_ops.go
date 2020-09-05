package core

// Pitched creates a new Note with a pitch by a (positive or negative) number of semi tones
func (n Note) Pitched(howManySemitones int) Note {
	if howManySemitones == 0 {
		return n
	}
	if n.IsRest() || n.IsPedalUp() || n.IsPedalDown() || n.IsPedalUpDown() {
		return n
	}
	simple := MIDItoNote(1.0, n.MIDI()+howManySemitones, 1.0)
	nn, _ := NewNote(simple.Name, simple.Octave, n.duration, simple.Accidental, n.Dotted, n.Velocity)
	return nn
}

func (n Note) Octaved(howmuch int) Note {
	if howmuch == 0 {
		return n
	}
	nn, _ := NewNote(n.Name, n.Octave+howmuch, n.duration, n.Accidental, n.Dotted, n.Velocity)
	return nn
}
