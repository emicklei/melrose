package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Group struct {
	Target []core.Sequenceable
}

func (p Group) S() core.Sequence {
	n := []core.Note{}
	for _, each := range p.Target {
		each.S().NotesDo(func(each core.Note) {
			n = append(n, each)
		})
	}
	return core.Sequence{Notes: [][]core.Note{n}}
}

func (p Group) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "group(")
	core.AppendStorexList(&b, true, p.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (p Group) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	return Group{Target: replacedAll(p.Target, from, to)}
}
