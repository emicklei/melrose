package img

import (
	"fmt"

	"github.com/fogleman/gg"
)

type PianoView struct {
	Low, High int // midi nr
}

func (p PianoView) DrawOn(gc *gg.Context, vp ViewPort) {
	yscale := vp.Height() / (float64(p.High+1) - float64(p.Low))
	margin := float64(gc.Height()) - vp.Top
	for nr := p.High; nr >= p.Low; nr-- {
		y := float64(nr-p.Low)*yscale + margin
		fmt.Print("y=", y, "h=", yscale)
		gc.DrawRectangle(vp.Left, y, vp.Width(), yscale)
		gc.Stroke()
	}
}
