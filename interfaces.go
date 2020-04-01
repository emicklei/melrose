package melrose

import "time"

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
	Record(deviceID int, stopAfterInactivity time.Duration) (Sequence, error)

	SetBeatsPerMinute(bpm float64)
	BeatsPerMinute() float64

	Close()
}
