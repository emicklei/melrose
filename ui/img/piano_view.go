package img

import (
	"image/color"

	"github.com/emicklei/melrose/core"
	"github.com/fogleman/gg"
)

type PianoView struct {
	Low, High int // midi nr
}

func (p PianoView) DrawOn(gc *gg.Context, vp ViewPort) {
	yscale := vp.Height() / (float64(p.High+1) - float64(p.Low))
	margin := float64(gc.Height()) - vp.Top
	for nr := p.High; nr >= p.Low; nr-- {
		isBlack := core.IsBlackKey(nr)
		y := float64(nr-p.Low)*yscale + margin
		if isBlack {
			gc.SetColor(color.Black)
			insetDrawRectangle(gc, vp.Left, y, vp.Width()*2/3, yscale, 2)
			gc.Fill()
		} else {
			// white
			gc.SetColor(color.White)
			gc.DrawRectangle(vp.Left, y, vp.Width(), yscale)
			gc.Fill()
			aboveIsBlack := core.IsBlackKey(nr + 1)
			if aboveIsBlack {
				gc.DrawRectangle(vp.Left+vp.Width()*2/3, y-(yscale/2), vp.Width()/3, yscale/2)
				gc.Fill()
			}
			underIsBlack := core.IsBlackKey(nr - 1)
			if underIsBlack {
				gc.DrawRectangle(vp.Left+vp.Width()*2/3, y+yscale, vp.Width()/3, yscale/2)
				gc.Fill()
			}
		}
	}
}

func insetDrawRectangle(gc *gg.Context, x, y, w, h, inset float64) {
	hi := inset / 2.0
	gc.DrawRectangle(x+hi, y+hi, w-inset, h-inset)
}
