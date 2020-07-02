package core

import (
	"fmt"
	"time"
)

type NoteChange struct {
	isOn     bool
	note     int64
	velocity int64
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

type Recording struct {
	timeline     *Timeline
	baseVelocity float32
}

func NewRecording() *Recording {
	tim := NewTimeline()
	return &Recording{
		timeline:     tim,
		baseVelocity: 70.0,
	}
}

func (r *Recording) Add(e NoteChange, when time.Time) {
	r.timeline.Schedule(e, when)
}

type noteChangeEvent struct {
	change NoteChange
	when   time.Time
}

func (r *Recording) String() string {
	return fmt.Sprintf("recording, #notes:%d", r.timeline.Len())
}

// Sequence is an alias for S()
func (r *Recording) Sequence() Sequence { return r.S() }

func (r *Recording) S() Sequence {
	activeNotes := map[int64]noteChangeEvent{}
	notes := []Note{}
	r.timeline.EventsDo(func(event TimelineEvent, when time.Time) {
		change := event.(NoteChange)
		if change.isOn {
			_, ok := activeNotes[change.note]
			if ok {
				// note was on ?
			} else {
				// new
				activeNotes[change.note] = noteChangeEvent{change: change, when: when}
			}
		} else {
			// note off
			_, ok := activeNotes[change.note]
			if !ok {
				// note was never on ?
			} else {
				//duration := when.Sub(active.when)
				note := MIDItoNote(int(change.note), int(change.velocity))
				notes = append(notes, note)
				delete(activeNotes, change.note)
			}
		}

	})
	return BuildSequence(notes)
}
