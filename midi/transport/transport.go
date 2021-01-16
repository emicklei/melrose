package transport

import "github.com/emicklei/melrose/core"

var Factory = func() Transporter { return nil }

type Transporter interface {
	HasInputCapability() bool
	PrintInfo(inID, outID int)
	DefaultOutputDeviceID() int
	NewMIDIOut(id int) (MIDIOut, error)
	NewMIDIIn(id int) (MIDIIn, error)
	Terminate()
	NewMIDIListener(MIDIIn) MIDIListener
}

type MIDIOut interface {
	WriteShort(status int64, data1 int64, data2 int64) error
	Close() error
	Abort() error
}

type MIDIIn interface {
	Close() error
}

type MIDIListener interface {
	Add(core.NoteListener)
	Remove(core.NoteListener)
	OnKey(core.Note, core.NoteListener)
	Start()
	Stop()
}
