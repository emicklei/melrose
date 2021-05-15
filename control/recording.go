package control

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Recording struct {
	deviceID     int
	timeline     *core.Timeline
	variableName string
}

func NewRecording(deviceID int, variableName string) *Recording {
	tim := core.NewTimeline()
	return &Recording{
		deviceID:     deviceID,
		timeline:     tim,
		variableName: variableName,
	}
}

type noteChangeEvent struct {
	change core.NoteChange
	when   time.Time
}

func (r *Recording) GetTargetFrom(other *Recording) {
	// only overwrite variable
	// listener may have been started so timeline is not empty, so device is listened to
	r.variableName = other.variableName
}

func (r *Recording) Play(ctx core.Context, at time.Time) error {
	ctx.Device().Listen(r.deviceID, r, true)
	return nil
}

func (r *Recording) Stop(ctx core.Context) error {
	// nothing there or already stopped
	if r.timeline.Len() == 0 {
		return nil
	}
	seq := r.S()
	if core.IsDebug() {
		notify.Debugf("recording.stop seq:%v", seq)
	}
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

func (r *Recording) S() core.Sequenceable {
	return buildExpressionFromTimeline(r.timeline)
}

func (r *Recording) NoteOn(channel int, n core.Note) {
	if core.IsDebug() {
		notify.Debugf("recording.noteon note:%v", n)
	}
	change := core.NewNoteChange(true, int64(n.MIDI()), int64(n.Velocity))
	r.timeline.Schedule(change, time.Now())
}

func (r *Recording) NoteOff(channel int, n core.Note) {
	if core.IsDebug() {
		notify.Debugf("recording.noteoff note:%v", n)
	}
	change := core.NewNoteChange(false, int64(n.MIDI()), int64(n.Velocity))
	r.timeline.Schedule(change, time.Now())
}

// ControlChange is ignored
func (r *Recording) ControlChange(channel, number, value int) {}

func buildExpressionFromTimeline(t *core.Timeline) core.Sequenceable {
	activeNotes := map[int]noteChangeEvent{}
	notes := []core.Note{}
	t.EventsDo(func(event core.TimelineEvent, when time.Time) {
		change := event.(core.NoteChange)
		if change.IsOn() {
			_, ok := activeNotes[change.Number()]
			if ok {
				// note was on ?
			} else {
				// new
				activeNotes[change.Number()] = noteChangeEvent{change: change, when: when}
			}
		} else {
			// note off
			hit, ok := activeNotes[change.Number()]
			if !ok {
				// note was never on ?
			} else {
				// when.Sub(active.when)fraction
				note, err := core.MIDItoNote(0.25, hit.change.Number(), int(hit.change.Velocity()))
				if err != nil {
					notify.Warnf("core.MIDItoNote error:%v", err)
				} else {
					notes = append(notes, note)
					delete(activeNotes, change.Number())
				}
			}
		}

	})
	return core.BuildSequence(notes)
}
