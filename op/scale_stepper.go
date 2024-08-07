package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type ScaleStepper struct {
	Scale  core.HasValue
	Target core.Sequenceable
	Count  core.HasValue
}

func (p ScaleStepper) S() core.Sequence {
	s := p.Scale.Value().(core.Scale)
	count := p.Count.Value().(int)
	notes := [][]core.Note{}
	for _, each := range p.Target.S().Notes {
		pair := []core.Note{}
		for _, other := range each {
			i := s.IndexOfNote(other)
			j := s.NoteAtIndex(i + count)
			pair = append(pair, j)
		}
		notes = append(notes, pair)
	}
	return core.Sequence{Notes: notes}
}

func (s ScaleStepper) Storex() string {
	return fmt.Sprintf("scale_stepper(%s,%s,%s)", core.Storex(s.Scale), core.Storex(s.Count), core.Storex(s.Target))
}

// Replaced is part of Replaceable
func (s ScaleStepper) Replaced(from, to core.Sequenceable) core.Sequenceable {
	return s
}
