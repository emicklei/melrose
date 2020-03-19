package melrose

import "fmt"

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

func (c Chord) String() string {
	return fmt.Sprintf("%v %s %s", c.start, minorMajor(c.minorOrMajor), inversion(c.inversion))
}

func minorMajor(m int) string {
	if m == Minor {
		return "minor"
	}
	if m == Major {
		return "major"
	}
	return "?"
}

func inversion(i int) string {
	if i == Ground {
		return "ground"
	}
	if i == Inversion1 {
		return "1st inversion"
	}
	if i == Inversion2 {
		return "2nd inversion"
	}
	if i == Inversion3 {
		return "3nd inversion"
	}
	return "?"
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
		inversion:      c.inversion,
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
