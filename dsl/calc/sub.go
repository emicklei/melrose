package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Sub struct {
	Left  any
	Right any
}

func (s Sub) Storex() string {
	return fmt.Sprintf("%s - %s", core.Storex(s.Left), core.Storex(s.Right))
}

func (s Sub) Value() any {
	isFloatOp := false
	if _, ok := s.Left.(float64); ok {
		isFloatOp = true
	}
	if _, ok := s.Right.(float64); ok {
		isFloatOp = true
	}

	if isFloatOp {
		l, _ := resolveFloatWithInt(s.Left)
		r, _ := resolveFloatWithInt(s.Right)
		return l - r
	}
	// integer op
	l, ok := resolveInt(s.Left)
	if !ok {
		l = 0
	}
	r, ok := resolveInt(s.Right)
	if !ok {
		r = 0
	}
	return l - r
}
