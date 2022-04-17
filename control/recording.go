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
	bpm          float64
}

func NewRecording(deviceID int, variableName string, bpm float64) *Recording {
	tim := core.NewTimeline()
	return &Recording{
		deviceID:     deviceID,
		timeline:     tim,
		variableName: variableName,
		bpm:          bpm,
	}
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

// Stop is part of Stoppable
func (r *Recording) Stop(ctx core.Context) error {
	// nothing there or already stopped
	if r.timeline.Len() == 0 {
		if core.IsDebug() {
			notify.Debugf("empty timeline on stop recording")
		}
		return nil
	}
	seq := r.S()
	if core.IsDebug() {
		notify.Debugf("recording.stop seq:%v", seq)
		// TODO temporary store the recording for the test in ui/img/draw_test.go
		fmt.Println("writing /tmp/melrose-recording.json")
		core.NotesEventsToFile(r.timeline.NoteEvents(), "/tmp/melrose-recording.json")
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
	periods := r.timeline.BuildNotePeriods()
	builder := core.NewSequenceBuilder(periods, r.bpm)
	return builder.Build()
}

func (r *Recording) NoteOn(channel int, n core.Note) {
	when := time.Now()
	change := core.NewNoteChange(true, int64(n.MIDI()), int64(n.Velocity))
	if core.IsDebug() {
		notify.Debugf("recording.noteon note:%v t:%s", n, when.Format("04:05.000"))
	}
	r.timeline.Schedule(change, when)
}

func (r *Recording) NoteOff(channel int, n core.Note) {
	when := time.Now()
	change := core.NewNoteChange(false, int64(n.MIDI()), int64(n.Velocity))
	if core.IsDebug() {
		notify.Debugf("recording.noteoff note:%v t:%s", n, when.Format("04:05.000"))
	}
	r.timeline.Schedule(change, when)
}

// ControlChange is ignored
func (r *Recording) ControlChange(channel, number, value int) {}
