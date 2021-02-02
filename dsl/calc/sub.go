package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Sub struct {
	Left  interface{}
	Right interface{}
}

func (s Sub) Storex() string {
	return fmt.Sprintf("%s - %s", core.Storex(s.Left), core.Storex(s.Right))
}

func (s Sub) Value() interface{} {
	l, ok := resolveInt(s.Left)
	if !ok {
		// try floats
		f, ok := s.floatValue()
		if ok {
			return f
		}
		l = 0
	}
	r, ok := resolveInt(s.Right)
	if !ok {
		r = 0
	}
	return l - r
}

func (s Sub) floatValue() (float64, bool) {
	l, ok := resolveFloat(s.Left)
	if !ok {
		return 0.0, false
	}
	r, ok := resolveFloat(s.Right)
	if !ok {
		return 0.0, false
	}
	return l - r, true
}
