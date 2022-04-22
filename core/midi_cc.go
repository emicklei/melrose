package core

type MIDIControlChange struct {
	ctx       Context
	parameter int
	value     HasValue
}

func NewMIDICC(p int, v HasValue) MIDIControlChange {
	return MIDIControlChange{
		parameter: p,
		value:     v,
	}
}

// S has the side effect of sending a MIDI CC message to the current output device
func (cc MIDIControlChange) S() Sequence {
	return EmptySequence
}

type SliderPosition struct {
	fraction float32
	dotted   bool
	value    int
}
