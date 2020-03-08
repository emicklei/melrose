package melrose

type Chord struct {
	minorOrMajor   int
	start          Note
	triadOrSeventh int
	inversion      int
}

func (n Note) Chord(modifiers ...int) Chord {
	zero := Chord{
		minorOrMajor:   Major,
		start:          n,
		triadOrSeventh: Triad,
		inversion:      Ground,
	}
	return zero.Modified(modifiers...)
}

func (c Chord) Modified(modifiers ...int) Chord {
	modified := c
	for _, each := range modifiers {
		switch each {
		case Major:
			modified.minorOrMajor = Major
		case Minor:
			modified.minorOrMajor = Minor
		case Ground:
			modified.inversion = Ground
		case Inversion1:
			modified.inversion = Inversion1
		case Inversion2:
			modified.inversion = Inversion2
		}
	}
	return modified
}

func (c Chord) Octaved(howMuch int) Chord {
	return Chord{
		minorOrMajor:   Minor,
		start:          c.start.Octaved(howMuch),
		triadOrSeventh: c.triadOrSeventh,
	}
}

func (c Chord) S() Sequence {
	notes := []Note{c.start}
	var semitones []int
	if Major == c.minorOrMajor {
		semitones = []int{4, 7}
	} else if Minor == c.minorOrMajor {
		semitones = []int{3, 7}
	}
	for _, each := range semitones {
		next := c.start.Pitched(each)
		notes = append(notes, next)
	}
	return Sequence{[][]Note{notes}}
}

func (c Chord) Join(j ...Joinable) Sequence {
	return c.S().Join(j...)
}

func (c Chord) NoteJoin(n Note) Sequence {
	return c.S().NoteJoin(n)
}

func (c Chord) SequenceJoin(s Sequence) Sequence {
	return c.S().SequenceJoin(s)
}
