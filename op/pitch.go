package op

import (
	"fmt"

	. "github.com/emicklei/melrose"
)

type Pitch struct {
	Target    Sequenceable
	Semitones Valueable
}

func (p Pitch) S() Sequence {
	return p.Target.S().Pitched(Int(p.Semitones))
}

func (p Pitch) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("pitch(%v,%s)", p.Semitones, s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (p Pitch) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(p, from) {
		return to
	}
	if IsIdenticalTo(p.Target, from) {
		return Pitch{Target: to, Semitones: p.Semitones}
	}
	// https://play.golang.org/p/qHbbK_sTo84
	if r, ok := p.Target.(Replaceable); ok {
		return r.Replaced(from, to)
	}
	return p
}
