package img

import (
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
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
	gc := gg.NewContext(1000, 400)

	vp := NewViewPort(10, 390, 990, 10)
	evts := tl.NoteEvents()
	nv := NotesView{Events: evts}
	nv.DrawOn(gc, vp)
	gc.SavePNG("out.png")
}

func TestPianoRoll(t *testing.T) {
	gc := gg.NewContext(1000, 400)
	gc.SetRGB(1.0, 1.0, 1.0)
	p := PianoView{Low: 30, High: 62}
	vp := NewViewPort(10, 390, 100, 10)
	p.DrawOn(gc, vp)
	gc.SetRGB(1.0, 0, 0)
	gc.DrawString("C#", 500, 200)
	gc.SavePNG("out.png")
}

func TestAxis(t *testing.T) {
	gc := gg.NewContext(100, 100)
	gc.SetRGB(1.0, 1.0, 1.0)
	gc.DrawRectangle(10, 50, 20, 40)
	gc.Fill()
	gc.SavePNG("out.png")
}
