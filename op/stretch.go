package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Stretch struct {
	target []core.Sequenceable
	factor core.HasValue
}

func NewStretch(factor core.HasValue, target []core.Sequenceable) Stretch {
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
	return Stretch{target: replacedAll(s.target, from, to), factor: s.factor}
}
