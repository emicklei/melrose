package op

import (
	"fmt"

	"github.com/emicklei/melrose"
)

type AtIndex struct {
	Target melrose.Sequenceable
	Index  melrose.Valueable
}

func (a AtIndex) S() melrose.Sequence {
	s := a.Target.S()
	i := melrose.Int(a.Index)
	if i < 1 {
		return melrose.Sequence{}
	}
	if i > len(s.Notes) {
		return melrose.Sequence{}
	}
	return melrose.BuildSequence(s.At(i - 1))
}

func (a AtIndex) Storex() string {
	if s, ok := a.Target.(melrose.Storable); ok {
		return fmt.Sprintf("at(%v,%s)", a.Index, s.Storex())
	}
	return ""
}

func NewAtIndex(index melrose.Valueable, target melrose.Sequenceable) AtIndex {
	return AtIndex{Target: target, Index: index}
}
