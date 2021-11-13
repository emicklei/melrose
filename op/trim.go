package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Trim struct {
	Start  core.HasValue
	End    core.HasValue
	Target core.Sequenceable
}

func (t Trim) S() core.Sequence {
	start, ok := t.Start.Value().(int)
	if !ok || start < 0 {
		start = 0
	}
	end, ok := t.End.Value().(int)
	if !ok || end < 0 {
		end = 0
	}
	notes := t.Target.S().Notes
	if end >= len(notes) {
		return core.EmptySequence
	}
	return core.Sequence{
		Notes: notes[start : len(notes)-end],
	}
}

func (t Trim) Storex() string {
	if s, ok := t.Target.(core.Storable); ok {
		return fmt.Sprintf("trim(%s,%s,%s)",
			core.Storex(t.Start),
			core.Storex(t.End),
			s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (t Trim) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(t, from) {
		return to
	}
	if core.IsIdenticalTo(t.Target, from) {
		return Trim{Start: t.Start, End: t.End, Target: to}
	}
	if rep, ok := t.Target.(core.Replaceable); ok {
		return Trim{Start: t.Start, End: t.End, Target: rep.Replaced(from, to)}
	}
	return t
}
