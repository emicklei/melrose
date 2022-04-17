package core

import "time"

type MIDIControlChange struct {
	parameter HasValue
	value     HasValue
}

func NewMIDICC(p, v HasValue) MIDIControlChange {
	return MIDIControlChange{
		parameter: p,
		value:     v,
	}
}

func (cc MIDIControlChange) Play(ctx Context, at time.Time) error {
	return nil
}
