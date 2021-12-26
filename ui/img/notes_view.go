package img

import (
	"github.com/emicklei/melrose/core"
	"github.com/fogleman/gg"
)

type NotesView struct {
	Events []core.NoteEvent
	BPM    float64
	// TODO BIAB
}

// gc 0,0 is top-left
func (v NotesView) DrawOn(gc *gg.Context) {
	if len(v.Events) == 0 || v.BPM == 0 {
		return
	}
	stats := core.NoteStatistics(v.Events)
	bottom := float64(gc.Height())
	xscale := float64(gc.Width()) / float64(stats.End.Sub(stats.Start).Milliseconds())
	yscale := float64(gc.Height()) / float64(stats.Highest-stats.Lowest+1)

	quarter := (core.WholeNoteDuration(v.BPM) / 4.0).Milliseconds()

	bar := 0
	for x := 0.0; x < float64(gc.Width()); x += float64(quarter) * xscale {
		bar++
		if bar == 4 {
			gc.SetRGB(200/256.0, 0.0, 0.0) // redisch
			bar = 0
		} else {
			gc.SetRGB(200/256.0, 200/256.0, 200/256.0) // grayish
		}
		gc.DrawLine(x, 0, x, float64(gc.Height()))
		gc.Stroke()
	}

	gc.SetRGB(62/256.0, 161/256.0, 11/256.0) // greenish
	for _, each := range v.Events {
		xs := float64(each.Start.Sub(stats.Start).Milliseconds()) * xscale
		xe := float64(each.End.Sub(stats.Start).Milliseconds()) * xscale

		gc.DrawRectangle(xs, bottom-(float64(each.Number-stats.Lowest+1)*yscale), xe-xs, yscale)
		gc.Fill()
	}
}
