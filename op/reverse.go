package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Reverse struct {
	Target []core.Sequenceable
}

func (r Reverse) S() core.Sequence {
	if len(r.Target) == 0 {
		return core.EmptySequence
	}
	return r.Target[0].S().Reversed()
}

func (r Reverse) Storex() string {
	if len(r.Target) == 0 {
		return ""
	}
	if s, ok := r.Target[0].(core.Storable); ok {
		return fmt.Sprintf("reverse(%s)", s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (r Reverse) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(r, from) {
		return to
	}
	return Reverse{Target: replacedAll(r.Target, from, to)}
}
