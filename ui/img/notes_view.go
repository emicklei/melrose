package img

import (
	"github.com/emicklei/melrose/core"
	"github.com/fogleman/gg"
)

type NotesView struct {
	Events []core.NoteEvent
}

func (v NotesView) DrawOn(gc *gg.Context, vp ViewPort) {
	if len(v.Events) == 0 {
		return
	}
	stats := core.NoteStatistics(v.Events)
	bottom := float64(gc.Height()) - vp.Bottom
	xscale := vp.Width() / float64(stats.End.Sub(stats.Start).Milliseconds())
	yscale := vp.Height() / float64(stats.Highest-stats.Lowest+1)

	gc.SetRGB(62/256.0, 161/256.0, 11/256.0)
	for _, each := range v.Events {
		xs := float64(each.Start.Sub(stats.Start).Milliseconds()) * xscale
		xe := float64(each.End.Sub(stats.Start).Milliseconds()) * xscale

		gc.DrawRectangle(vp.Left+xs, bottom-(float64(each.Number-stats.Lowest+1)*yscale), xe-xs, yscale)
		gc.Fill()
	}
}
