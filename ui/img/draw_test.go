package img

import (
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/op"
	"github.com/fogleman/gg"
)

func TestNotesViewWidth(t *testing.T) {
	bpm := 120.0
	start := time.Now()
	tests := []struct {
		name     string
		duration time.Duration
		want     int
	}{
		{name: "thirty-second note", duration: core.WholeNoteDuration(bpm) / 32, want: 2},
		{name: "sixteenth note", duration: core.WholeNoteDuration(bpm) / 16, want: 4},
		{name: "whole note", duration: core.WholeNoteDuration(bpm), want: 64},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			view := NotesView{
				Events: []core.NoteEvent{{Start: start, End: start.Add(test.duration), Number: 60}},
				BPM:    bpm,
			}
			if got := view.Width(); got != test.want {
				t.Fatalf("Width() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestNotesViewHeight(t *testing.T) {
	view := NotesView{Events: []core.NoteEvent{
		{Number: 60},
		{Number: 64},
	}}
	if got, want := view.Height(), 20; got != want {
		t.Fatalf("Height() = %d, want %d", got, want)
	}
}

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
	evts := tl.NoteEvents()
	nv := NotesView{Events: evts, BPM: 10.0}
	gc := gg.NewContext(nv.Width(), nv.Height())
	nv.DrawOn(gc)
	gc.SavePNG("TestDraw.png")
}

func TestRecordedTimeline(t *testing.T) {
	t.Skip("TODO put recorded file in testdata")
	bpm := 120.0
	// TODO stored from control/recording.go:54
	events := core.NoteEventsFromFile("/tmp/melrose-recording.json")
	t.Log("event count:", len(events))
	nv := NotesView{Events: events, BPM: bpm}
	gc := gg.NewContext(nv.Width(), nv.Height())
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
		nv := NotesView{Events: tim.NoteEvents(), BPM: bpm}
		gc := gg.NewContext(nv.Width(), nv.Height())
		nv.DrawOn(gc)
		gc.SavePNG("TestRecorded_PROCESSED.png")
	}
}

func TestScaleInputSequenceBuilder(t *testing.T) {
	bpm := 120.0
	s1, _ := core.NewScale("8C")
	s2, _ := core.NewScale("8C3")
	seq := op.Merge{
		Target: []core.Sequenceable{s1, s2},
	}
	tim := core.NewTimeline()
	d := midi.NewOutputDevice(0, nil, 0, tim)
	d.Play(core.NoCondition, seq, bpm, time.Now())
	nv := NotesView{Events: tim.NoteEvents(), BPM: bpm}
	gc := gg.NewContext(nv.Width(), nv.Height())
	nv.DrawOn(gc)
	gc.SavePNG("TestRecorded_SCALE.png")
}
