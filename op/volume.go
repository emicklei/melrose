package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Volume struct {
	Target []core.Sequenceable
	Value  core.HasValue
}

func (v Volume) S() core.Sequence {
	actual := core.ValueOf(v.Value)
	volume, ok := actual.(int)
	if !ok || volume < 0 || volume > 127 {
		notify.Warnf("[op.Volume] volume must be an integer in [0..127], got (%T) %v", actual, actual)
		return core.EmptySequence
	}

	sequence := Join{Target: v.Target}.S()
	for groupIndex, group := range sequence.Notes {
		for noteIndex, note := range group {
			sequence.Notes[groupIndex][noteIndex] = note.WithVelocity(volume)
		}
	}
	return sequence
}

func (v Volume) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "volume(%s", core.Storex(v.Value))
	core.AppendStorexList(&b, false, v.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (v Volume) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(v, from) {
		return to
	}
	return Volume{Target: replacedAll(v.Target, from, to), Value: v.Value}
}
