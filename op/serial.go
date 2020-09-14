package op

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emicklei/melrose/core"
)

type Serial struct {
	Target []core.Sequenceable
}

func (a Serial) S() core.Sequence {
	n := []core.Note{}
	for _, each := range a.Target {
		each.S().NotesDo(func(each core.Note) {
			n = append(n, each)
		})
	}
	return core.BuildSequence(n)
}

// Storex is part of Storable
func (a Serial) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "ungroup(")
	for i, each := range a.Target {
		s, ok := each.(core.Storable)
		if !ok {
			return ""
		}
		fmt.Fprintf(&b, "%s", s.Storex())
		if i < len(a.Target)-1 {
			io.WriteString(&b, ",")
		}
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (a Serial) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(a, from) {
		return to
	}
	return Serial{Target: replacedAll(a.Target, from, to)}
}
