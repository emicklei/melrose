package op

import (
	"fmt"
	"github.com/emicklei/melrose/core"
)

type Parallel struct {
	Target core.Sequenceable
}

func (p Parallel) S() core.Sequence {
	n := []core.Note{}
	p.Target.S().NotesDo(func(each core.Note) {
		n = append(n, each)
	})
	return core.Sequence{Notes: [][]core.Note{n}}
}

func (p Parallel) Storex() string {
	if s, ok := p.Target.(core.Storable); ok {
		return fmt.Sprintf("parallel(%s)", s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (p Parallel) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	if core.IsIdenticalTo(p.Target, from) {
		return Parallel{Target: to}
	}
	if rep, ok := p.Target.(core.Replaceable); ok {
		return Parallel{Target: rep.Replaced(from, to)}
	}
	return p
}
