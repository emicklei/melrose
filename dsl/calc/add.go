package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Add struct {
	Left  interface{}
	Right interface{}
}

func (a Add) Storex() string {
	return fmt.Sprintf("%s + %s", core.Storex(a.Left), core.Storex(a.Right))
}

func (a Add) Value() interface{} {
	l, ok := resolveInt(a.Left)
	if !ok {
		// try floats
		f, ok := a.floatValue()
		if ok {
			return f
		}
		l = 0
	}
	r, ok := resolveInt(a.Right)
	if !ok {
		r = 0
	}
	return l + r
}

func (a Add) floatValue() (float64, bool) {
	l, ok := resolveFloat(a.Left)
	if !ok {
		return 0.0, false
	}
	r, ok := resolveFloat(a.Right)
	if !ok {
		return 0.0, false
	}
	return l + r, true
}
