package core

// Pitched creates a new Note with a pitch by a (positive or negative) number of semi tones
func (n Note) Pitched(howManySemitones int) Note {
	if howManySemitones == 0 {
		return n
	}
	if n.IsRest() || n.IsPedalUp() || n.IsPedalDown() || n.IsPedalUpDown() {
		return n
	}
	simple, err := MIDItoNote(1.0, n.MIDI()+howManySemitones, 1.0)
	if err != nil {
		panic(err)
	}
	p := MakeNote(simple.Name, simple.Octave, n.fraction, simple.Accidental, n.Dotted, n.Velocity)
	// handle tied notes
	for _, each := range n.tied {
		p = p.WithTiedNote(each.Pitched(howManySemitones))
	}
	return p
}

func (n Note) Octaved(howmuch int) Note {
	if howmuch == 0 {
		return n
	}
	p := MakeNote(n.Name, n.Octave+howmuch, n.fraction, n.Accidental, n.Dotted, n.Velocity)
	// handle tied notes
	for _, each := range n.tied {
		p = p.WithTiedNote(each.Octaved(howmuch))
	}
	return p
}

func (n Note) Stretched(f float32) Note {
	if f == 1.0 {
		return n
	}
	p := MakeNote(n.Name, n.Octave, n.fraction*f, n.Accidental, n.Dotted, n.Velocity)
	// handle tied notes
	for _, each := range n.tied {
		p = p.WithTiedNote(each.Stretched(f))
	}
	return p
}
