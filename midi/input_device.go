package midi

type InputDevice struct {
	stream         MIDIIn
	defaultChannel int
}
