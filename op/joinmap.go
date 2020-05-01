package op

import (
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

type JoinMapper struct {
	Target  melrose.Valueable
	Indices [][]int
}

func (j JoinMapper) S() melrose.Sequence {
	join, ok := j.Target.Value().(Join)
	if !ok {
		return melrose.Sequence{}
	}
	source := join.Target
	target := []melrose.Sequenceable{}
	for i, indexEntry := range j.Indices {
		mapped := []melrose.Sequenceable{}
		for j, each := range indexEntry {
			if each < 0 || each > len(source) {
				notify.Print(notify.Warningf("index out of join range: [%d][%d]=%d, len=%d", i+1, j+1, each, len(source)))
			} else {
				mapped = append(mapped, source[each-1])
			}
		}
		target = append(target, Join{Target: mapped})
	}
	return Join{Target: target}.S()
}

func NewJoinMapper(v melrose.Valueable, indices string) JoinMapper {
	return JoinMapper{Target: v, Indices: parseIndices(indices)}
}
