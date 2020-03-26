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
	PrintInfo()
	Play(seq Sequence, echo bool)
	SetBeatsPerMinute(bpm float64)
	BeatsPerMinute() float64
	Close()
}
