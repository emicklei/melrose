package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type JoinMap struct {
	target  core.HasValue
	pattern core.HasValue
}

func (j JoinMap) Storex() string {
	return fmt.Sprintf("joinmap('%s',%s)", formatIndices(j.indices()), core.Storex(j.target))
}

func (j JoinMap) indices() [][]int {
	ps := core.String(j.pattern)
	return parseIndices(ps)
}

func (j JoinMap) S() core.Sequence {
	join, ok := j.target.Value().(Join)
	if !ok {
		return core.EmptySequence
	}
	source := join.Target
	target := []core.Sequenceable{}
	for i, indexGroup := range j.indices() {
		if len(indexGroup) == 1 {
			// single
			if j.check(i, 0, indexGroup[0], len(source)) {
				target = append(target, source[indexGroup[0]-1])
			} else {
				target = append(target, core.Rest4) // TODO what should be the duration?
			}
		} else {
			// group
			notes := []core.Note{}
			for g, each := range indexGroup {
				if j.check(i, g, each, len(source)) {
					notes = append(notes, source[each-1].S().Notes[0]...)
				} else {
					target = append(target, core.Rest4) // TODO what should be the duration?
				}
			}
			target = append(target, Group{Target: core.BuildSequence(notes)})
		}
	}
	return Join{Target: target}.S()
}

func (j JoinMap) check(index, subindex, value, length int) bool { // indices are zero-based
	if value < 1 || value > length {
		notify.Warnf("index out of join range: [%d][%d]=%d, len=%d, using a rest(=) instead", index+1, subindex+1, value, length)
		return false
	}
	return true
}

func NewJoinMap(v core.HasValue, indices core.HasValue) JoinMap {
	return JoinMap{target: v, pattern: indices}
}

// Replaced is part of Replaceable
func (j JoinMap) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(j, from) {
		return to
	}
	join, ok := j.target.Value().(Join)
	if !ok {
		return j
	}
	return JoinMap{target: core.On(join.Replaced(from, to)), pattern: j.pattern}
}
