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

	// Play schedules all the notes on the timeline using a BPM (beats-per-minute).
	// Returns the end time of the last played Note.
	Play(seq Sequenceable, bpm float64) time.Time

	Record(deviceID int, stopAfterInactivity time.Duration) (Sequence, error)

	Timeline() *Timeline

	SetBaseVelocity(velocity int)
	SetEchoNotes(echo bool)

	Close()
}

type LoopController interface {
	Start()
	Stop()

	SetBPM(bpm float64)
	BPM() float64

	SetBIAB(biab int)
	BIAB() int

	Begin(l *Loop)
	End(l *Loop)
}

type MapFunc func(seq Sequenceable) Sequenceable

// TODO experiment
type Mappeable interface {
	Map(m MapFunc) Mappeable
}
