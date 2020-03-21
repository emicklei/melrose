package melrose

type Transformer interface {
	Transform(Sequence) Sequence
}

type Sequenceable interface {
	Storable
	S() Sequence
}

type Storable interface {
	Storex() string
}

type AudioDevice interface {
	Play(s Sequence)
	BeatsPerMinute(bpm float64)
	Close()
}
