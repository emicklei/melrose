package core

import (
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

type Sequenceable interface {
	S() Sequence
}

type NoteConvertable interface {
	ToNote() (Note, error)
}

type Storable interface {
	Storex() string
}

type Indexable interface {
	At(i int) Sequenceable
}

type Nextable interface {
	Next() interface{}
}

type AudioDevice interface {
	DefaultDeviceIDs() (inputDeviceID, outputDeviceID int)

	// Per device specific commands
	Command(args []string) notify.Message

	// Handle generic setting
	HandleSetting(name string, values []interface{}) error

	// Play schedules all the notes on the timeline using a BPM (beats-per-minute).
	// Returns the end time of the last played Note.
	Play(condition Condition, seq Sequenceable, bpm float64, beginAt time.Time) (endingAt time.Time)

	HasInputCapability() bool
	Listen(deviceID int, who NoteListener, startOrStop bool)

	// if a key is pressed on a device then play or stop a function
	// if fun is nil then uninstall the binding
	OnKey(ctx Context, deviceID int, channel int, note Note, fun HasValue) error

	// Schedule put an event on the timeline at a begin
	Schedule(event TimelineEvent, beginAt time.Time)

	// Record(ctx Context) (*Recording, error)
	Reset()
	Close() error
}

type LoopController interface {
	Start()
	Stop()
	Reset()

	SetBPM(bpm float64)
	BPM() float64

	SetBIAB(biab int)
	BIAB() int

	BeatsAndBars() (int64, int64)
	Plan(bars int64, seq Sequenceable)

	SettingNotifier(handler func(control LoopController))
}

type Replaceable interface {
	// Returns a new value in which any occurrences of "from" are replaced by "to".
	Replaced(from, to Sequenceable) Sequenceable
}

type HasValue interface {
	Value() interface{}
}

type Inspectable interface {
	Inspect(i Inspection)
}

type Playable interface {
	Play(ctx Context, at time.Time) error
}

type Stoppable interface {
	Stop(ctx Context) error
	IsPlaying() bool
}

type VariableStorage interface {
	NameFor(value interface{}) string
	Get(key string) (interface{}, bool)
	Put(key string, value interface{})
	Delete(key string)
	Variables() map[string]interface{}
}

type Context interface {
	Control() LoopController
	Device() AudioDevice
	Variables() VariableStorage
	Environment() *sync.Map
	WithCondition(c Condition) Context
	Capabilities() *Capabilities
}

// WorkingDirectory is a key in a context environment.
const WorkingDirectory = "shell.pwd"
const EditorLineStart = "editor.line.start"
const EditorLineEnd = "editor.line.end"

// TODO makue users use Play with a Context that can have a Condition
type Evaluatable interface {
	Evaluate(ctx Context) error
}

type NoteListener interface {
	NoteOn(channel int, note Note)
	NoteOff(channel int, note Note)
	ControlChange(channel, number, value int)
}

type Conditional interface {
	Condition() Condition
}

type Condition func() bool

var (
	NoCondition   Condition = nil
	TrueCondition Condition = func() bool { return true }
)

type NameAware interface {
	VariableName(yours string)
}

type HasIndex interface {
	Index() HasValue
}
