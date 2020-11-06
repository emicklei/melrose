package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Pitch struct {
	Target    core.Sequenceable
	Semitones core.Valueable
}

func (p Pitch) S() core.Sequence {
	return p.Target.S().Pitched(core.Int(p.Semitones))
}

func (p Pitch) Storex() string {
	return fmt.Sprintf("pitch(%v,%s)", p.Semitones, core.Storex(p.Target))
}

// Replaced is part of Replaceable
func (p Pitch) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	if core.IsIdenticalTo(p.Target, from) {
		return Pitch{Target: to, Semitones: p.Semitones}
	}
	// https://play.golang.org/p/qHbbK_sTo84
	if r, ok := p.Target.(core.Replaceable); ok {
		return r.Replaced(from, to)
	}
	return p
}
