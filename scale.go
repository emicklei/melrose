package melrose

type Scale struct {
	minorOrMajor int
	start        Note
}

func (n Note) Scale(modifiers ...int) Scale {
	d := Scale{
		start:        n,
		minorOrMajor: Major,
	}
	return d.Modified(modifiers...)
}

func (s Scale) Modified(modifiers ...int) Scale {
	modified := s
	for _, each := range modifiers {
		switch each {
		case Major:
			modified.minorOrMajor = Major
		case Minor:
			modified.minorOrMajor = Minor
		}
	}
	return modified
}

func (s Scale) Storex() string {
	// TODO handle minor major
	return s.start.Storex() + ".Scale()"
}

func (s Scale) S() Sequence {
	notes := []Note{s.start}
	var semitones []int
	if Major == s.minorOrMajor {
		semitones = []int{2, 2, 1, 2, 2, 2}
	} else if Minor == s.minorOrMajor {
		semitones = []int{2, 1, 2, 2, 1, 2}
	}
	next := s.start
	for _, each := range semitones {
		next = next.Pitched(each)
		notes = append(notes, next)
	}
	return BuildSequence(notes)
}
