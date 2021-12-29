package img

import (
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/op"
	"github.com/fogleman/gg"
)

func sampleTimeline() *core.Timeline {
	tim := core.NewTimeline()
	now := time.Now()
	tim.Schedule(core.NewNoteChange(true, 21, 1), now.Add(0*time.Second))
	tim.Schedule(core.NewNoteChange(false, 21, 1), now.Add(1*time.Second))
	tim.Schedule(core.NewNoteChange(true, 21, 1), now.Add(3*time.Second))
	tim.Schedule(core.NewNoteChange(true, 22, 1), now.Add(4*time.Second))
	tim.Schedule(core.NewNoteChange(false, 22, 1), now.Add(5*time.Second))
	tim.Schedule(core.NewNoteChange(false, 21, 1), now.Add(6*time.Second))

	tim.Schedule(core.NewNoteChange(true, 21, 1), now.Add(8*time.Second))
	tim.Schedule(core.NewNoteChange(true, 23, 1), now.Add(9*time.Second))
	tim.Schedule(core.NewNoteChange(false, 21, 1), now.Add(10*time.Second))
	tim.Schedule(core.NewNoteChange(false, 23, 1), now.Add(11*time.Second))

	tim.Schedule(core.NewNoteChange(true, 21, 1), now.Add(13*time.Second))
	tim.Schedule(core.NewNoteChange(true, 24, 1), now.Add(14*time.Second))
	tim.Schedule(core.NewNoteChange(false, 21, 1), now.Add(15*time.Second))
	tim.Schedule(core.NewNoteChange(false, 24, 1), now.Add(15*time.Second))

	tim.Schedule(core.NewNoteChange(true, 21, 1), now.Add(17*time.Second))
	tim.Schedule(core.NewNoteChange(true, 25, 1), now.Add(17*time.Second))
	tim.Schedule(core.NewNoteChange(false, 25, 1), now.Add(18*time.Second))
	tim.Schedule(core.NewNoteChange(false, 21, 1), now.Add(19*time.Second))

	return tim.ZeroStarting()
}

func TestDraw(t *testing.T) {
	tl := sampleTimeline()
	gc := gg.NewContext(1000, 150)

	evts := tl.NoteEvents()
	nv := NotesView{Events: evts, BPM: 10.0}
	nv.DrawOn(gc)
	gc.SavePNG("TestDraw.png")
}

func TestRecordedTimeline(t *testing.T) {
	bpm := 120.0
	events := core.NoteEventsFromFile("/tmp/melrose-recording.json")
	gc := gg.NewContext(500, 50)
	nv := NotesView{Events: events, BPM: bpm}
	nv.DrawOn(gc)
	gc.SavePNG("TestRecorded_RAW.png")

	{
		periods := core.NoteEventsToPeriods(events)
		b := core.NewSequenceBuilder(periods, bpm)
		seq := b.Build()
		t.Log(seq)
		tim := core.NewTimeline()
		d := midi.NewOutputDevice(0, nil, 0, tim)
		d.Play(core.NoCondition, seq, bpm, time.Now())
		gc := gg.NewContext(500, 50)
		nv := NotesView{Events: tim.NoteEvents(), BPM: bpm}
		nv.DrawOn(gc)
		gc.SavePNG("TestRecorded_PROCESSED.png")
	}
}

func TestScaleInputSequenceBuilder(t *testing.T) {
	bpm := 120.0
	s1, _ := core.NewScale(2, "8C")
	s2, _ := core.NewScale(2, "8C3")
	seq := op.Merge{
		Target: []core.Sequenceable{s1, s2},
	}
	tim := core.NewTimeline()
	d := midi.NewOutputDevice(0, nil, 0, tim)
	d.Play(core.NoCondition, seq, bpm, time.Now())
	gc := gg.NewContext(500, 100)
	nv := NotesView{Events: tim.NoteEvents(), BPM: bpm}
	nv.DrawOn(gc)
	gc.SavePNG("TestRecorded_SCALE.png")
}
