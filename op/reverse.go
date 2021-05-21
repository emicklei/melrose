package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Reverse struct {
	Target core.Sequenceable
}

func (r Reverse) S() core.Sequence {
	return r.Target.S().Reversed()
}

func (r Reverse) Storex() string {
	if s, ok := r.Target.(core.Storable); ok {
		return fmt.Sprintf("reverse(%s)", s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (r Reverse) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(r, from) {
		return to
	}
	if core.IsIdenticalTo(r.Target, from) {
		return Reverse{Target: to}
	}
	if tr, ok := r.Target.(core.Replaceable); ok {
		return Reverse{Target: tr.Replaced(from, to)}
	}
	return r
}
