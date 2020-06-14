package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose"
)

type Join struct {
	Target []melrose.Sequenceable
}

func (j Join) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "join(")
	appendStorexList(&b, true, j.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (j Join) S() melrose.Sequence {
	if len(j.Target) == 0 {
		return melrose.EmptySequence
	}
	head := j.Target[0].S()
	for i := 1; i < len(j.Target); i++ {
		head = head.SequenceJoin(j.Target[i].S())
	}
	return head
}

// Replaced is part of Replaceable
func (j Join) Replaced(from, to melrose.Sequenceable) melrose.Sequenceable {
	if melrose.IsIdenticalTo(j, from) {
		return to
	}
	return Join{Target: replacedAll(j.Target, from, to)}
}
