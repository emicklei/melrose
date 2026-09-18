package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Transpose struct {
	Target    []core.Sequenceable
	Semitones core.HasValue
}

func (p Transpose) S() core.Sequence {
	if len(p.Target) == 0 {
		return core.EmptySequence
	}
	return p.Target[0].S().Pitched(core.Int(p.Semitones))
}

func (p Transpose) Storex() string {
	if len(p.Target) == 0 {
		return fmt.Sprintf("transpose(%s,nil)", core.Storex(p.Semitones))
	}
	return fmt.Sprintf("transpose(%s,%s)", core.Storex(p.Semitones), core.Storex(p.Target[0]))
}

// Replaced is part of Replaceable
func (p Transpose) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	return Transpose{Target: replacedAll(p.Target, from, to), Semitones: p.Semitones}
}
