package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type NumberCompare struct {
	Left     interface{}
	Right    interface{}
	Operator string
}

func (a NumberCompare) Storex() string {
	return fmt.Sprintf("%s %s %s", core.Storex(a.Left), a.Operator, core.Storex(a.Right))
}

func (a NumberCompare) Value() interface{} {
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
	switch a.Operator {
	case "<":
		return l < r
	case "<=":
		return l <= r
	case ">":
		return l > r
	case ">=":
		return l >= r
	case "!=":
		return l != r
	case "==":
		return l == r
	default:
		return false
	}
}

// floatValue operate on left and right as floats
func (a NumberCompare) floatValue() (bool, bool) {
	l, ok := resolveFloat(a.Left)
	if !ok {
		return false, false
	}
	r, ok := resolveFloat(a.Right)
	if !ok {
		return false, false
	}
	switch a.Operator {
	case "<":
		return l < r, true
	case "<=":
		return l <= r, true
	case ">":
		return l > r, true
	case ">=":
		return l >= r, true
	case "!=":
		return l != r, true
	case "==":
		return l == r, true
	default:
		return false, false
	}
}
