package img

import (
	"fmt"
	"image/color"

	"github.com/emicklei/melrose/core"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
)

type NotesView struct {
	Events []core.NoteEvent
	BPM    float64
	// TODO BIAB
}

const (
	bottomPadding = 20
)

// gc 0,0 is top-left
func (v NotesView) DrawOn(gc draw2d.GraphicContext, width, height int) {
	if len(v.Events) == 0 || v.BPM == 0 {
		return
	}
	stats := core.NoteStatistics(v.Events)
	viewHeight := float64(height - bottomPadding)
	bottom := viewHeight
	if stats.End.Sub(stats.Start).Milliseconds() == 0 {
		return
	}
	xscale := float64(width) / float64(stats.End.Sub(stats.Start).Milliseconds())
	yscale := viewHeight / float64(stats.Highest-stats.Lowest+1)

	quarterMS := (core.WholeNoteDuration(v.BPM) / 4.0).Milliseconds()

	bar := 0
	for x := 0.0; x < float64(width); x += float64(quarterMS) * xscale {
		isNewBar := bar%4 == 0
		if isNewBar {
			gc.SetStrokeColor(color.RGBA{R: 200, G: 0, B: 0, A: 255}) // redisch
		} else {
			gc.SetStrokeColor(color.RGBA{R: 200, G: 200, B: 200, A: 255}) // grayish
		}
		gc.MoveTo(x, 0)
		gc.LineTo(x, viewHeight)
		gc.Stroke()

		// draw time
		if isNewBar {
			gc.SetFillColor(color.RGBA{R: 0, G: 0, B: 0, A: 255}) // black
			momentS := (x / xscale) / 1000.0
			gc.FillStringAt(fmt.Sprintf("%.1fs", momentS), x, viewHeight+15)
		}
		bar++
	}

	gc.SetFillColor(color.RGBA{R: 62, G: 161, B: 11, A: 255}) // greenish
	for _, each := range v.Events {
		xs := float64(each.Start.Sub(stats.Start).Milliseconds()) * xscale
		xe := float64(each.End.Sub(stats.Start).Milliseconds()) * xscale
		ys := bottom - (float64(each.Number-stats.Lowest+1) * yscale)

		draw2dkit.Rectangle(gc, xs, ys, xe, ys+yscale)
		gc.Fill()

		// draw note label
		note, err := core.MIDItoNote(0.25, each.Number, each.Velocity)
		if err == nil {
			gc.SetFillColor(color.RGBA{R: 0, G: 0, B: 0, A: 255}) // black
			gc.FillStringAt(note.String(), xs, ys-2)
			gc.SetFillColor(color.RGBA{R: 62, G: 161, B: 11, A: 255}) // back to greenish
		}
	}
}
