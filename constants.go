package melrose

import "regexp"

const (
	Sharp = iota
	Flat
	Eight
	Quarter
	Half
	Whole
	Dot

	Left
	Right

	PrintAsSpecified

	Major
	Minor // Natural
	HarmonicMinor
	MelodicMinor
)
const (
	// NonRestNoteNames maps a tone to an index (C=0)
	NonRestNoteNames = "CDEFGAB"
)

var (
	noteRegexp = regexp.MustCompile("([½¼⅛248]?)([CDEFGABr])([#♯_♭]?)(\\.?)([0-9]?)")

	// noteMidiOffsets maps a tone index (C=0) to the number of semitones on the scale
	noteMidiOffsets = []int{0, 2, 4, 5, 7, 9, 11}
)
