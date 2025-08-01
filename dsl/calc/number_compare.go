package calc

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type NumberCompare struct {
	Left     any
	Right    any
	Operator string
}

func (a NumberCompare) Storex() string {
	return fmt.Sprintf("%s %s %s", core.Storex(a.Left), a.Operator, core.Storex(a.Right))
}

func (a NumberCompare) Value() any {
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

	l, ok := resolveInt(a.Left)
	if !ok {
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
