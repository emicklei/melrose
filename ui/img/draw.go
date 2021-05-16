package img

import (
	"github.com/emicklei/melrose/core"
	"github.com/fogleman/gg"
)

func Draw(ncl []core.NoteEvent) {
	if len(ncl) == 0 {
		return
	}
	dc := gg.NewContext(1000, 400)
	dc.SetRGB(1.0, 1.0, 1.0)

	bottom := float64(dc.Height())

	stats := core.NoteStatistics(ncl)
	//fmt.Println(stats)

	xscale := float64(dc.Width()) / float64(stats.End.Sub(stats.Start).Milliseconds())
	yscale := float64(dc.Height())/float64(stats.Highest) - float64(stats.Lowest)

	for _, each := range ncl {
		xs := float64(each.Start.Sub(stats.Start).Milliseconds()) * xscale
		xe := float64(each.End.Sub(stats.Start).Milliseconds()) * xscale
		dc.DrawRectangle(xs, bottom-float64(each.Number)*yscale, xe-xs, yscale)
		dc.Fill()
	}
	dc.SavePNG("out.png")
}
