package melrose

type Joinable interface {
	JoinNote(Note) Sequence
	JoinSequence(Sequence) Sequence
}

type Joiner interface {
	Join(j ...Joinable) Sequence
}
