package midi

import (
	"github.com/emicklei/melrose/midi/transport"
)

type InputDevice struct {
	id       int
	echo     bool
	listener transport.MIDIListener
}

func NewInputDevice(id int, in transport.MIDIIn, t transport.Transporter) *InputDevice {
	return &InputDevice{
		id:       id,
		echo:     false,
		listener: t.NewMIDIListener(in), /// newListener(in.(*portmidi.Stream)), // TODO
	}
}

func (i *InputDevice) stopListener() {
	i.listener.Stop()
}
