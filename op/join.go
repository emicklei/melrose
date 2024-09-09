package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Join struct {
	Target []core.Sequenceable
}

// SequenceableList is part of core.HasSequenceables
func (j Join) Sequenceables() []core.Sequenceable {
	return j.Target
}

func (j Join) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "join(")
	core.AppendStorexList(&b, true, j.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (j Join) S() core.Sequence {
	if len(j.Target) == 0 {
		return core.EmptySequence
	}
	joined := j.Target[0].S()
	for i := 1; i < len(j.Target); i++ {
		joined = joined.SequenceJoin(j.Target[i].S())
	}
	return joined
}

// Replaced is part of Replaceable
func (j Join) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(j, from) {
		return to
	}
	return Join{Target: replacedAll(j.Target, from, to)}
}
