package midi

import "github.com/rakyll/portmidi"

type InputDevice struct {
	id       int
	stream   MIDIIn
	echo     bool
	listener *listener
}

func NewInputDevice(id int, in MIDIIn) *InputDevice {
	return &InputDevice{
		id:       id,
		stream:   in,
		echo:     false,
		listener: newListener(in.(*portmidi.Stream)), // TODO
	}
}

func (i *InputDevice) stopListener() {
	i.listener.stop()
}
