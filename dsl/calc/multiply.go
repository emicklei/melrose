package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Multiply struct {
	Left  interface{}
	Right interface{}
}

func (m Multiply) Storex() string {
	return fmt.Sprintf("%s * %s", core.Storex(m.Left), core.Storex(m.Right))
}

func (m Multiply) Value() interface{} {
	l, ok := resolveInt(m.Left)
	if !ok {
		// try floats
		f, ok := m.floatValue()
		if ok {
			return f
		}
		l = 0
	}
	r, ok := resolveInt(m.Right)
	if !ok {
		r = 0
	}
	return l * r
}

func (m Multiply) floatValue() (float64, bool) {
	l, ok := resolveFloat(m.Left)
	if !ok {
		return 0.0, false
	}
	r, ok := resolveFloat(m.Right)
	if !ok {
		return 0.0, false
	}
	return l * r, true
}
