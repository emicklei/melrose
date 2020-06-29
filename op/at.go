package op

import (
	"fmt"
	"github.com/emicklei/melrose/core"
)

type AtIndex struct {
	Target core.Sequenceable
	Index  core.Valueable
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
	if s, ok := a.Target.(core.Storable); ok {
		return fmt.Sprintf("at(%v,%s)", a.Index, s.Storex())
	}
	return ""
}

func NewAtIndex(index core.Valueable, target core.Sequenceable) AtIndex {
	return AtIndex{Target: target, Index: index}
}
