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
	return fmt.Sprintf("%s + %s", core.Storex(m.Left), core.Storex(m.Right))
}

func (m Multiply) Value() interface{} {
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
