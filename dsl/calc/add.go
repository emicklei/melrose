package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/op"
)

type Add struct {
	Left  interface{}
	Right interface{}
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
