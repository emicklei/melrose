package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Repeat struct {
	Target []core.Sequenceable
	Times  core.Valueable
}

func (r Repeat) S() core.Sequence {
	times := core.Int(r.Times)
	repeated := []core.Sequenceable{}
	for i := 0; i < times; i++ {
		for _, each := range r.Target {
			repeated = append(repeated, each.S())
		}
	}
	return Join{Target: repeated}.S()
}

func (r Repeat) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "repeat(%s", core.Storex(r.Times))
	appendStorexList(&b, false, r.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (r Repeat) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(r, from) {
		return to
	}
	return Repeat{Target: replacedAll(r.Target, from, to), Times: r.Times}
}
