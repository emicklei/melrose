package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Transpose struct {
	Target    core.Sequenceable
	Semitones core.HasValue
}

func (p Transpose) S() core.Sequence {
	return p.Target.S().Pitched(core.Int(p.Semitones))
}

func (p Transpose) Storex() string {
	return fmt.Sprintf("transpose(%s,%s)", core.Storex(p.Semitones), core.Storex(p.Target))
}

// Replaced is part of Replaceable
func (p Transpose) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	if core.IsIdenticalTo(p.Target, from) {
		return Transpose{Target: to, Semitones: p.Semitones}
	}
	// https://play.golang.org/p/qHbbK_sTo84
	if r, ok := p.Target.(core.Replaceable); ok {
		return r.Replaced(from, to)
	}
	return p
}
