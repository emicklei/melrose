package img

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
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
	gc := gg.NewContext(500, 50)

	evts := tl.NoteEvents()
	nv := NotesView{Events: evts}
	nv.DrawOn(gc)
	gc.SavePNG("TestDraw.png")
}

func TestScale(t *testing.T) {
	tim := core.NewTimeline()
	d := midi.NewOutputDevice(0, nil, 0, tim)
	s, _ := core.NewScale(2, "C")
	d.Play(core.NoCondition, s, 120.0, time.Now())

	gc := gg.NewContext(500, 50)

	evts := tim.NoteEvents()
	nv := NotesView{Events: evts}
	nv.DrawOn(gc)
	gc.SavePNG("TestScale.png")
}

func TestRecordedTimeline(t *testing.T) {
	in, _ := os.Open("rec.json")
	defer in.Close()
	dec := json.NewDecoder(in)
	list := []core.StorableNoteChange{}
	if err := dec.Decode(&list); err != nil {
		t.Fatal(err)
	}
	// for _, each := range list {
	// 	fmt.Printf("%#v\n", each)
	// }

	gc := gg.NewContext(500, 50)
	events := core.ConvertToNoteEvents(list)
	// for _, each := range events {
	// 	fmt.Printf("%#v\n", each)
	// }

	nv := NotesView{Events: events}
	nv.DrawOn(gc)
	gc.SavePNG("TestRecordedTimeline.png")
}
