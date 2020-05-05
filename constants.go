package melrose

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

// https://nl.wikipedia.org/wiki/Dynamiek_(muziek)
const (
	Pianississimo = 16 // ----
	Pianissimo    = 33 // ---
	Piano         = 49 // --
	MezzoPiano    = 64 // -
	Normal        = 72
	MezzoForte    = 80  // +
	Forte         = 96  // ++
	Fortissimo    = 112 // +++
	Fortississimo = 127 // ++++
)
