package control

import (
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
)

func TestRecordingStop(t *testing.T) {
	var r interface{} = new(Recording)
	_, ok := r.(core.Stoppable)
	if !ok {
		t.Fail()
	}
}

func sampleRecording() *Recording {
	tim := core.NewTimeline()
	now := time.Now()
	tim.Schedule(core.NewNoteChange(true, 1, 1), now.Add(0*time.Second))
	tim.Schedule(core.NewNoteChange(false, 1, 1), now.Add(1*time.Second))
	tim.Schedule(core.NewNoteChange(true, 1, 1), now.Add(3*time.Second))
	tim.Schedule(core.NewNoteChange(true, 2, 1), now.Add(4*time.Second))
	tim.Schedule(core.NewNoteChange(false, 2, 1), now.Add(5*time.Second))
	tim.Schedule(core.NewNoteChange(false, 1, 1), now.Add(6*time.Second))

	tim.Schedule(core.NewNoteChange(true, 1, 1), now.Add(8*time.Second))
	tim.Schedule(core.NewNoteChange(true, 3, 1), now.Add(9*time.Second))
	tim.Schedule(core.NewNoteChange(false, 1, 1), now.Add(10*time.Second))
	tim.Schedule(core.NewNoteChange(false, 3, 1), now.Add(11*time.Second))

	tim.Schedule(core.NewNoteChange(true, 1, 1), now.Add(13*time.Second))
	tim.Schedule(core.NewNoteChange(true, 4, 1), now.Add(14*time.Second))
	tim.Schedule(core.NewNoteChange(false, 1, 1), now.Add(15*time.Second))
	tim.Schedule(core.NewNoteChange(false, 4, 1), now.Add(15*time.Second))

	tim.Schedule(core.NewNoteChange(true, 1, 1), now.Add(17*time.Second))
	tim.Schedule(core.NewNoteChange(true, 5, 1), now.Add(17*time.Second))
	tim.Schedule(core.NewNoteChange(false, 5, 1), now.Add(18*time.Second))
	tim.Schedule(core.NewNoteChange(false, 1, 1), now.Add(19*time.Second))

	return &Recording{timeline: tim.ZeroStarting()}
}

func TestRecordingSequence(t *testing.T) {
	rec := sampleRecording()
	t.Log(rec.timeline.Len())
}
