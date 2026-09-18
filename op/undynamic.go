package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Undynamic struct {
	Target []core.Sequenceable
}

func (u Undynamic) S() core.Sequence {
	if len(u.Target) == 0 {
		return core.EmptySequence
	}
	n := []core.Note{}
	u.Target[0].S().NotesDo(func(each core.Note) {
		each.Velocity = core.Normal
		n = append(n, each)
	})
	return core.BuildSequence(n)
}

func (u Undynamic) Storex() string {
	if len(u.Target) == 0 {
		return ""
	}
	if s, ok := u.Target[0].(core.Storable); ok {
		return fmt.Sprintf("undynamic(%s)", s.Storex())
	}
	return ""
}
