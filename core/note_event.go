package core

import "time"

type NoteEvent struct {
	Start, End time.Time
	Number     int
	Velocity   int
}

func (e NoteEvent) WithEnd(end time.Time) NoteEvent {
	return NoteEvent{Start: e.Start, Number: e.Number, Velocity: e.Velocity, End: end}
}

type NoteChange struct {
	isOn     bool
	note     int64
	velocity int64
}

func (n NoteChange) Number() int {
	return int(n.note)
}

func (n NoteChange) Velocity() int {
	return int(n.velocity)
}

func (n NoteChange) IsOn() bool {
	return n.isOn
}

func NewNoteChange(isOn bool, midiNr int64, velocity int64) NoteChange {
	return NoteChange{
		isOn:     isOn,
		note:     midiNr,
		velocity: velocity,
	}
}

func (n NoteChange) Handle(tim *Timeline, when time.Time) {
	// NOP
}
