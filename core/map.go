package core

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/notify"
)

type Map struct {
	Target      HasValue
	Replaceable HasValue
	Each        Sequenceable
}

func (c Map) S() Sequence {
	return SequenceableList{Target: c.Sequenceables()}.S()
}

func (c Map) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "map(%s,%s)", Storex(c.Target), Storex(c.Replaceable))
	return b.String()
}

func (c Map) Sequenceables() []Sequenceable {
	tv := c.Target.Value()
	t, ok := tv.(HasSequenceables)
	if !ok {
		notify.Warnf("target does not have sequences")
		return []Sequenceable{}
	}
	rv := c.Replaceable.Value()
	r, ok := rv.(Replaceable)
	if !ok {
		notify.Warnf("function does not allow replacement")
		return []Sequenceable{}
	}
	targets := make([]Sequenceable, len(t.Sequenceables()))
	for i, each := range t.Sequenceables() {
		targets[i] = r.Replaced(c.Each, each)
	}
	return targets
}

func (c Map) Replaced(from, to Sequenceable) Sequenceable {
	if from == Sequenceable(c) {
		return to
	}
	newTarget := c.Target
	if t, ok := c.Target.Value().(Replaceable); ok {
		newTarget = On(t.Replaced(from, to))
	}
	newReplaceable := c.Replaceable
	if r, ok := c.Replaceable.Value().(Replaceable); ok {
		newReplaceable = On(r.Replaced(from, to))
	}
	return Map{Target: newTarget, Replaceable: newReplaceable, Each: c.Each}
}
