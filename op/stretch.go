package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Stretch struct {
	target []core.Sequenceable
	factor core.Valueable
}

func NewStretch(factor core.Valueable, target []core.Sequenceable) Stretch {
	return Stretch{
		target: target,
		factor: factor,
	}
}

func (s Stretch) S() core.Sequence {
	return Join{Target: s.target}.S().Stretched(core.Float(s.factor))
}

// Storex is part of Storable
func (s Stretch) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "stretch(%s", core.Storex(s.factor))
	core.AppendStorexList(&b, false, s.target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (s Stretch) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(s, from) {
		return to
	}
	// TODO
	// if core.IsIdenticalTo(s.factor, from) {
	// 	return Stretch{target: s.target, factor: from}
	// }
	return Stretch{target: replacedAll(s.target, from, to)}
}
