package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Bare struct {
	Target []core.Sequenceable
}

// TODO remove pedals, bindings
func (b Bare) S() core.Sequence {
	notes := [][]core.Note{}
	for _, each := range b.Target {
		simple := each.S().NoFractions().NoDynamics().NoRests()
		notes = append(notes, simple.Notes...)
	}
	return core.Sequence{Notes: notes}
}

func (b Bare) Storex() string {
	var bb bytes.Buffer
	fmt.Fprintf(&bb, "bare(")
	core.AppendStorexList(&bb, true, b.Target)
	fmt.Fprintf(&bb, ")")
	return bb.String()
}

// Replaced is part of Replaceable
func (b Bare) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(b, from) {
		return to
	}
	return Bare{Target: replacedAll(b.Target, from, to)}
}
