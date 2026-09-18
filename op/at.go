package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type AtIndex struct {
	Target []core.Sequenceable
	Index  core.HasValue
}

func (a AtIndex) S() core.Sequence {
	if len(a.Target) == 0 {
		return core.EmptySequence
	}
	s := a.Target[0].S()
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
	if len(a.Target) == 0 {
		return fmt.Sprintf("at(%v,nil)", core.Storex(a.Index))
	}
	return fmt.Sprintf("at(%v,%s)", core.Storex(a.Index), core.Storex(a.Target[0]))
}

func NewAtIndex(index core.HasValue, target core.Sequenceable) AtIndex {
	return AtIndex{Target: []core.Sequenceable{target}, Index: index}
}
