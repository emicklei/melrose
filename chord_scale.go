package melrose

// Scale returns a new Sequence of Notes starting from a Note in Major of Minor
func Scale(start Note, minorOrMajor int) Sequence {
	notes := []Note{start}
	var semitones []int
	if Major == minorOrMajor {
		semitones = []int{2, 2, 1, 2, 2, 2, 1}
	} else if Minor == minorOrMajor {
		semitones = []int{2, 1, 2, 2, 1, 2, 2}
	}
	next := start
	for _, each := range semitones {
		next = next.ModifiedPitch(each)
		notes = append(notes, next)
	}
	return BuildSequence(notes)
}

func Chord(start Note, minorOrMajor int) Sequence {
	notes := []Note{start}
	var semitones []int
	if Major == minorOrMajor {
		semitones = []int{4, 7}
	} else if Minor == minorOrMajor {
		semitones = []int{3, 7}
	}
	for _, each := range semitones {
		next := start.ModifiedPitch(each)
		notes = append(notes, next)
	}
	return Sequence{[][]Note{notes}}

}
