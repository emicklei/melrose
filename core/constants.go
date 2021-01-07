package core

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
	Dominant
	Augmented
	Diminished
	// https://en.wikipedia.org/wiki/Chord_(music)#Common_types_of_chords
	Triad
	Seventh
	Sixth

	Ground
	Inversion1
	Inversion2
	Inversion3
)

// TODO typed constants
type inversion int

type quality int

// const (
// 	Major = quality(0)
// 	Minor = quality(1)
// )

// https://www.cs.cmu.edu/~music/cmsip/readings/MIDI%20tutorial%20for%20programmers.html
const (
	VelocityPPPP = 8   // ----
	VelocityPPP  = 20  // ---
	VelocityPP   = 31  // --
	VelocityP    = 42  // -
	VelocityMP   = 53  // o-
	Normal       = 59  // o
	VelocityMF   = 64  // o+
	VelocityF    = 80  // +
	VelocityFF   = 96  // ++
	VelocityFFF  = 112 // +++
	VelocityFFFF = 127 // ++++
)
