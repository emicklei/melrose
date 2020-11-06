package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Group struct {
	Target core.Sequenceable
}

func (p Group) S() core.Sequence {
	n := []core.Note{}
	p.Target.S().NotesDo(func(each core.Note) {
		n = append(n, each)
	})
	return core.Sequence{Notes: [][]core.Note{n}}
}

func (p Group) Storex() string {
	return fmt.Sprintf("group(%s)", core.Storex(p.Target))
}

// Replaced is part of Replaceable
func (p Group) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	if core.IsIdenticalTo(p.Target, from) {
		return Group{Target: to}
	}
	if rep, ok := p.Target.(core.Replaceable); ok {
		return Group{Target: rep.Replaced(from, to)}
	}
	return p
}
