package img

import (
	"fmt"

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

	start := ncl[0].Start
	end := ncl[len(ncl)-1].End

	xscale := float64(1000.0) / float64(end.Sub(start).Milliseconds())

	for _, each := range ncl {
		fmt.Println(each)
		xs := float64(each.Start.Sub(start).Milliseconds()) * xscale
		xe := float64(each.End.Sub(start).Milliseconds()) * xscale
		dc.DrawRectangle(float64(xs), bottom-float64(each.Number)*100, xe-xs, 10)
		dc.Stroke()
	}
	dc.SavePNG("out.png")
}
