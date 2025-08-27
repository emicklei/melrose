package img

import (
	"image"
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/op"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dsvg"
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
	evts := tl.NoteEvents()
	nv := NotesView{Events: evts, BPM: 10.0}

	// PNG
	dest := image.NewRGBA(image.Rect(0, 0, 1000, 170))
	gc := draw2dimg.NewGraphicContext(dest)
	nv.DrawOn(gc, 1000, 170)
	draw2dimg.SaveToPngFile("TestDraw.png", dest)

	// SVG
	destSvg := draw2dsvg.NewSvg()
	gcSvg := draw2dsvg.NewGraphicContext(destSvg)
	nv.DrawOn(gcSvg, 1000, 170)
	draw2dsvg.SaveToSvgFile("TestDraw.svg", destSvg)
}

func TestRecordedTimeline(t *testing.T) {
	t.Skip("TODO put recorded file in testdata")
	bpm := 120.0
	// TODO stored from control/recording.go:54
	events := core.NoteEventsFromFile("/tmp/melrose-recording.json")
	t.Log("event count:", len(events))
	nv := NotesView{Events: events, BPM: bpm}

	dest := image.NewRGBA(image.Rect(0, 0, 1000, 120))
	gc := draw2dimg.NewGraphicContext(dest)
	nv.DrawOn(gc, 1000, 120)
	draw2dimg.SaveToPngFile("TestRecorded_RAW.png", dest)

	{
		periods := core.NoteEventsToPeriods(events)
		b := core.NewSequenceBuilder(periods, bpm)
		seq := b.Build()
		t.Log(seq)
		tim := core.NewTimeline()
		d := midi.NewOutputDevice(0, nil, 0, tim)
		d.Play(core.NoCondition, seq, bpm, time.Now())

		dest2 := image.NewRGBA(image.Rect(0, 0, 500, 70))
		gc2 := draw2dimg.NewGraphicContext(dest2)
		nv2 := NotesView{Events: tim.NoteEvents(), BPM: bpm}
		nv2.DrawOn(gc2, 500, 70)
		draw2dimg.SaveToPngFile("TestRecorded_PROCESSED.png", dest2)
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

	dest := image.NewRGBA(image.Rect(0, 0, 500, 120))
	gc := draw2dimg.NewGraphicContext(dest)
	nv := NotesView{Events: tim.NoteEvents(), BPM: bpm}
	nv.DrawOn(gc, 500, 120)
	draw2dimg.SaveToPngFile("TestRecorded_SCALE.png", dest)
}
