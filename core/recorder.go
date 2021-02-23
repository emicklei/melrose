package core

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/notify"
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
	deviceID     int
	timeline     *Timeline
	variableName string
}

func NewRecording(deviceID int, variableName string) *Recording {
	tim := NewTimeline()
	return &Recording{
		deviceID:     deviceID,
		timeline:     tim,
		variableName: variableName,
	}
}

type noteChangeEvent struct {
	change NoteChange
	when   time.Time
}

func (r *Recording) GetTargetFrom(other *Recording) {
	// only overwrite variable
	// listener may have been started so timeline is not empty, so device is listened to
	r.variableName = other.variableName
}

// Sequence is an alias for S()
func (r *Recording) Sequence() Sequence { return r.S() }

func (r *Recording) Play(ctx Context, at time.Time) error {
	ctx.Device().Listen(r.deviceID, r, true)
	return nil
}

func (r *Recording) Stop(ctx Context) error {
	seq := r.S()
	ctx.Variables().Put(r.variableName, seq)
	ctx.Device().Listen(r.deviceID, r, false)
	// flush
	r.timeline.Reset()
	return nil
}

func (r *Recording) IsPlaying() bool { return true }

func (r *Recording) Storex() string {
	return fmt.Sprintf("record(device(%d,%s))", r.deviceID, r.variableName)
}

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
				// when.Sub(active.when)fraction
				note, err := MIDItoNote(0.25, int(change.note), int(change.velocity))
				if err != nil {
					notify.Panic(err)
				} else {
					notes = append(notes, note)
					delete(activeNotes, change.note)
				}
			}
		}

	})
	return BuildSequence(notes)
}

func (r *Recording) NoteOn(n Note) {
	change := NewNoteChange(true, int64(n.MIDI()), int64(n.Velocity))
	r.timeline.Schedule(change, time.Now())
}

func (r *Recording) NoteOff(n Note) {
	change := NewNoteChange(false, int64(n.MIDI()), int64(n.Velocity))
	r.timeline.Schedule(change, time.Now())
}

// ControlChange is ignored
func (r *Recording) ControlChange(channel, number, value int) {}
