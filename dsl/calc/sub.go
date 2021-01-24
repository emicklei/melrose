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
		l = 0
	}
	r, ok := resolveInt(s.Right)
	if !ok {
		r = 0
	}
	return l - r
}
