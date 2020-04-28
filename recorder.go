package melrose

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

type Recorder struct {
	timeline     *Timeline
	baseVelocity float32
}

func NewRecorder() *Recorder {
	tim := NewTimeline()
	return &Recorder{
		timeline:     tim,
		baseVelocity: 70.0,
	}
}

func (r *Recorder) Add(e NoteChange, when time.Time) {
	r.timeline.Schedule(e, when)
}

type noteChangeEvent struct {
	change NoteChange
	when   time.Time
}

func (r *Recorder) BuildSequence() Sequence {
	activeNotes := map[int64]noteChangeEvent{}
	notes := []Note{}
	r.timeline.eventsDo(func(event TimelineEvent, when time.Time) {
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
			active, ok := activeNotes[change.note]
			if !ok {
				// note was never on ?
			} else {
				duration := when.Sub(active.when)
				fmt.Printf("%v\n", duration)
				note := MIDItoNote(int(change.note), 1.0)
				notes = append(notes, note)
				delete(activeNotes, change.note)
			}
		}

	})
	return BuildSequence(notes)
}
