package core

import (
	"fmt"
	"time"
)

type NoteEvent struct {
	Start, End time.Time
	Number     int
	Velocity   int
}

func (e NoteEvent) WithEnd(end time.Time) NoteEvent {
	return NoteEvent{Start: e.Start, Number: e.Number, Velocity: e.Velocity, End: end}
}

type NoteEventStatistics struct {
	Start, End      time.Time
	Lowest, Highest int
}

func NoteStatistics(list []NoteEvent) (stats NoteEventStatistics) {
	if len(list) == 0 {
		return stats
	}
	stats.Start = list[0].Start
	stats.End = list[len(list)-1].End
	stats.Lowest = 127
	for _, each := range list {
		if each.Number < stats.Lowest {
			stats.Lowest = each.Number
		}
		if each.Number > stats.Highest {
			stats.Highest = each.Number
		}
	}
	return
}

const eventTimeFormat = "03:04:05.000000000"

func (s NoteEventStatistics) String() string {
	return fmt.Sprintf("start=%s,end=%s,low=%d,high=%d", s.Start.Format(eventTimeFormat), s.End.Format(eventTimeFormat), s.Lowest, s.Highest)
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
