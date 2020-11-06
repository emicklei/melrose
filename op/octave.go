package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Octave struct {
	Target []core.Sequenceable
	Offset core.Valueable
}

func (o Octave) S() core.Sequence {
	return Join{Target: o.Target}.S().Octaved(core.Int(o.Offset))
}

func (o Octave) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "octave(%v", o.Offset)
	core.AppendStorexList(&b, false, o.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (o Octave) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(o, from) {
		return to
	}
	return Octave{Target: replacedAll(o.Target, from, to), Offset: o.Offset}
}
