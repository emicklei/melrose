package img

// ViewPort uses screen coordinates, Top is top of screen > Bottom.
type ViewPort struct {
	Left, Top, Right, Bottom float64
}

func NewViewPort(left, top, right, bottom int) ViewPort {
	return ViewPort{Left: float64(left), Top: float64(top), Right: float64(right), Bottom: float64(bottom)}
}

func (v ViewPort) Width() float64 {
	return float64(v.Right - v.Left)
}

func (v ViewPort) Height() float64 {
	return float64(v.Top - v.Bottom)
}
