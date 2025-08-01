package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/op"
)

type Add struct {
	Left  any
	Right any
}

func (a Add) S() core.Sequence {
	ls, ok := a.Left.(core.Sequenceable)
	if !ok {
		return core.EmptySequence
	}
	rs, ok := a.Right.(core.Sequenceable)
	if !ok {
		return core.EmptySequence
	}
	return op.Join{
		Target: []core.Sequenceable{ls, rs},
	}.S()
}

func (a Add) Storex() string {
	return fmt.Sprintf("%s + %s", core.Storex(a.Left), core.Storex(a.Right))
}

func (a Add) Value() any {
	isFloatOp := false
	if _, ok := a.Left.(float64); ok {
		isFloatOp = true
	}
	if _, ok := a.Right.(float64); ok {
		isFloatOp = true
	}

	if isFloatOp {
		l, _ := resolveFloatWithInt(a.Left)
		r, _ := resolveFloatWithInt(a.Right)
		return l + r
	}
	// integer op
	l, ok := resolveInt(a.Left)
	if !ok {
		l = 0
	}
	r, ok := resolveInt(a.Right)
	if !ok {
		r = 0
	}
	return l + r
}
