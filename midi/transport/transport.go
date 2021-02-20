package transport

import "github.com/emicklei/melrose/core"

// https://www.midi.org/specifications-old/item/table-1-summary-of-midi-message
const (
	noteOn        int64 = 0x90 // 10010000 , 144
	noteOff       int64 = 0x80 // 10000000 , 128
	controlChange int64 = 0xB0 // 10110000 , 176
	noteAllOff    int64 = 0x78 // 01111000 , 120  (not 123 because sustain)
	sustainPedal  int64 = 0x40
)

var Factory = func() Transporter { return nil }

var Initializer = func() {}

type Transporter interface {
	HasInputCapability() bool
	PrintInfo(inID, outID int)
	DefaultOutputDeviceID() int
	DefaultInputDeviceID() int
	NewMIDIOut(id int) (MIDIOut, error)
	NewMIDIIn(id int) (MIDIIn, error)
	NewMIDIListener(MIDIIn) MIDIListener
}

type MIDIOut interface {
	WriteShort(status int64, data1 int64, data2 int64) error
	Close() error
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
