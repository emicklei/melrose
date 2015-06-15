package melrose

type Joinable interface {
	// result is Note + Joinable
	NoteJoin(Note) Sequence
	// result is Sequence + Joinable
	SequenceJoin(Sequence) Sequence
}

type Joiner interface {
	Join(j ...Joinable) Sequence
}
