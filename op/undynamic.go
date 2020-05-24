package op

import (
	"fmt"

	. "github.com/emicklei/melrose"
)

type Undynamic struct {
	Target Sequenceable
}

func (u Undynamic) S() Sequence {
	n := []Note{}
	u.Target.S().NotesDo(func(each Note) {
		each.Velocity = Normal
		n = append(n, each)
	})
	return BuildSequence(n)
}

func (u Undynamic) Storex() string {
	if s, ok := u.Target.(Storable); ok {
		return fmt.Sprintf("undynamic(%s)", s.Storex())
	}
	return ""
}
