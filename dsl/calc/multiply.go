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
	isFloatOp := false
	if _, ok := m.Left.(float64); ok {
		isFloatOp = true
	}
	if _, ok := m.Right.(float64); ok {
		isFloatOp = true
	}

	if isFloatOp {
		l, _ := resolveFloatWithInt(m.Left)
		r, _ := resolveFloatWithInt(m.Right)
		return l * r
	}
	// integer op
	l, ok := resolveInt(m.Left)
	if !ok {
		l = 0
	}
	r, ok := resolveInt(m.Right)
	if !ok {
		r = 0
	}
	return l * r
}
