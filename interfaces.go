package melrose

import (
	"time"

	"github.com/emicklei/melrose/notify"
)

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
	// Per device specific commands
	Command(args []string) notify.Message

	Play(seq Sequenceable)
	Record(deviceID int, stopAfterInactivity time.Duration) (Sequence, error)

	SetBeatsPerMinute(bpm float64)
	SetBaseVelocity(velocity int)
	SetEchoNotes(echo bool)

	Close()
}
