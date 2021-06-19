package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type AtIndex struct {
	Target core.Sequenceable
	Index  core.HasValue
}

func (a AtIndex) S() core.Sequence {
	s := a.Target.S()
	i := core.Int(a.Index)
	if i < 1 {
		return core.EmptySequence
	}
	if i > len(s.Notes) {
		return core.EmptySequence
	}
	return core.BuildSequence(s.At(i - 1))
}

func (a AtIndex) Storex() string {
	return fmt.Sprintf("at(%v,%s)", core.Storex(a.Index), core.Storex(a.Target))
}

func NewAtIndex(index core.HasValue, target core.Sequenceable) AtIndex {
	return AtIndex{Target: target, Index: index}
}
