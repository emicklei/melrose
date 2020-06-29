package op

import (
	"fmt"
	"github.com/emicklei/melrose/core"
)

type Undynamic struct {
	Target core.Sequenceable
}

func (u Undynamic) S() core.Sequence {
	n := []core.Note{}
	u.Target.S().NotesDo(func(each core.Note) {
		each.Velocity = core.Normal
		n = append(n, each)
	})
	return core.BuildSequence(n)
}

func (u Undynamic) Storex() string {
	if s, ok := u.Target.(core.Storable); ok {
		return fmt.Sprintf("undynamic(%s)", s.Storex())
	}
	return ""
}
