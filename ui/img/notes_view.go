package img

import (
	"math"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/fogleman/gg"
)

const (
	pixelsPerWholeNote = 64
	pixelsPerSemitone  = 4
)

type NotesView struct {
	Events []core.NoteEvent
	BPM    float64
	// TODO BIAB
}

func (v NotesView) Width() int {
	if len(v.Events) == 0 || v.BPM <= 0 {
		return 1
	}
	stats := core.NoteStatistics(v.Events)
	return max(1, int(math.Ceil(v.durationToPixels(stats.End.Sub(stats.Start)))))
}

func (v NotesView) Height() int {
	if len(v.Events) == 0 {
		return 1
	}
	stats := core.NoteStatistics(v.Events)
	return (stats.Highest - stats.Lowest + 1) * pixelsPerSemitone
}

func (v NotesView) durationToPixels(duration time.Duration) float64 {
	return float64(duration) * pixelsPerWholeNote / float64(core.WholeNoteDuration(v.BPM))
}

// gc 0,0 is top-left
func (v NotesView) DrawOn(gc *gg.Context) {
	if len(v.Events) == 0 || v.BPM == 0 {
		return
	}
	stats := core.NoteStatistics(v.Events)
	bottom := float64(gc.Height())
	yscale := float64(gc.Height()) / float64(stats.Highest-stats.Lowest+1)

	bar := 0
	for x := 0.0; x < float64(gc.Width()); x += pixelsPerWholeNote / 4.0 {
		if bar == 4 {
			gc.SetRGB(200/256.0, 0.0, 0.0) // redish
			bar = 0
		} else {
			gc.SetRGB(200/256.0, 200/256.0, 200/256.0) // grayish
		}
		gc.DrawLine(x, 0, x, float64(gc.Height()))
		gc.Stroke()
		bar++
	}

	gc.SetRGB(62/256.0, 161/256.0, 11/256.0) // greenish
	for _, each := range v.Events {
		xs := v.durationToPixels(each.Start.Sub(stats.Start))
		xe := v.durationToPixels(each.End.Sub(stats.Start))

		gc.DrawRectangle(xs, bottom-(float64(each.Number-stats.Lowest+1)*yscale), xe-xs, yscale)
		gc.Fill()
	}
}
