package midi

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
		listener: nil,
	}
}
