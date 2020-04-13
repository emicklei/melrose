package melrose

import "time"

type Transformer interface {
	Transform(Sequence) Sequence
}

type Sequenceable interface {
	S() Sequence
}

type Storable interface {
	Storex() string
}

type AudioDevice interface {
	PrintInfo()

	Play(seq Sequenceable, echo bool)
	Record(deviceID int, stopAfterInactivity time.Duration) (Sequence, error)

	SetDefaultChannel(channel int)
	SetBeatsPerMinute(bpm float64)
	BeatsPerMinute() float64

	Close()
}
