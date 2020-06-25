package op

import (
	"fmt"

	"github.com/emicklei/melrose"
	. "github.com/emicklei/melrose"
)

type Parallel struct {
	Target Sequenceable
}

func (p Parallel) S() Sequence {
	n := []Note{}
	p.Target.S().NotesDo(func(each Note) {
		n = append(n, each)
	})
	return Sequence{Notes: [][]Note{n}}
}

func (p Parallel) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("parallel(%s)", s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (p Parallel) Replaced(from, to melrose.Sequenceable) melrose.Sequenceable {
	if melrose.IsIdenticalTo(p, from) {
		return to
	}
	if melrose.IsIdenticalTo(p.Target, from) {
		return Parallel{Target: to}
	}
	if rep, ok := p.Target.(melrose.Replaceable); ok {
		return Parallel{Target: rep.Replaced(from, to)}
	}
	return p
}
