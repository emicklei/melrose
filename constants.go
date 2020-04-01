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
	HarmonicMinor
	MelodicMinor

	// https://en.wikipedia.org/wiki/Chord_(music)#Common_types_of_chords
	MajorTriad
	MajorSeventh
	Triad
	Seventh

	Ground
	Inversion1
	Inversion2
	Inversion3

	// https://nl.wikipedia.org/wiki/Dynamiek_(muziek)
	Pianissimo // ---
	Piano      // --
	MezzoPiano // -
	MezzoForte // +
	Forte      // ++
	Fortissimo // +++
)
const (
	// NonRestNoteNames maps a tone to an index (C=0)
	NonRestNoteNames = "CDEFGAB"
)
