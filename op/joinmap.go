package op

import (
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

type JoinMapper struct {
	Target  melrose.Valueable
	Indices [][]int
}

func (j JoinMapper) Storex() string {
	return ""
}

func (j JoinMapper) S() melrose.Sequence {
	join, ok := j.Target.Value().(Join)
	if !ok {
		return melrose.Sequence{}
	}
	source := join.Target
	target := []melrose.Sequenceable{}
	for i, indexGroup := range j.Indices {
		if len(indexGroup) == 1 {
			// single
			if j.check(i, 0, indexGroup[0], len(source)) {
				target = append(target, source[indexGroup[0]-1])
			} else {
				target = append(target, melrose.Rest4) // what should be the duration?
			}
		} else {
			// group
			notes := []melrose.Note{}
			for g, each := range indexGroup {
				if j.check(i, g, each, len(source)) {
					notes = append(notes, source[each-1].S().Notes[0]...)
				} else {
					target = append(target, melrose.Rest4) // what should be the duration?
				}
			}
			target = append(target, Parallel{Target: melrose.BuildSequence(notes)})
		}
	}
	return Join{Target: target}.S()
}

func (j JoinMapper) check(index, subindex, value, length int) bool { // indices are zero-based
	if value < 1 || value > length {
		notify.Print(notify.Warningf("index out of join range: [%d][%d]=%d, len=%d, using a rest(=) instead", index+1, subindex+1, value, length))
		return false
	}
	return true
}

func NewJoinMapper(v melrose.Valueable, indices string) JoinMapper {
	return JoinMapper{Target: v, Indices: parseIndices(indices)}
}
