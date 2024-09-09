package core

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/notify"
)

type Collect struct {
	Target      HasValue
	Replaceable HasValue
	Each        Sequenceable
}

func (c Collect) S() Sequence {
	tv := c.Target.Value()
	t, ok := tv.(HasSequenceables)
	if !ok {
		notify.Warnf("target does not have sequences")
		return EmptySequence
	}
	rv := c.Replaceable.Value()
	r, ok := rv.(Replaceable)
	if !ok {
		notify.Warnf("function does not allow replacement")
		return EmptySequence
	}
	targets := make([]Sequenceable, len(t.Sequenceables()))
	for i, each := range t.Sequenceables() {
		targets[i] = r.Replaced(c.Each, each)
	}
	return SequenceableList{Target: targets}.S()
}

func (c Collect) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "collect(%s,%s)", Storex(c.Target), Storex(c.Replaceable))
	return b.String()
}
