package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
)

type SequenceMapper struct {
	Target  core.Sequenceable
	Indices [][]int
	Pattern core.Valueable
}

func (p SequenceMapper) S() core.Sequence {
	if p.Pattern == nil {
		return p.Target.S()
	}
	sPattern := core.String(p.Pattern)
	if len(sPattern) == 0 {
		return p.Target.S()
	}
	indices := parseIndices(sPattern)
	seq := p.Target.S()
	groups := [][]core.Note{}
	for _, indexEntry := range indices {
		mappedGroup := []core.Note{}
		for j, each := range indexEntry {
			if each < 1 || each > len(seq.Notes) {
				notify.Print(notify.Warningf("index out of sequence range: %d, len=%d", j+1, len(seq.Notes)))
			} else {
				// TODO what if sequence had a multi note group?
				mappedGroup = append(mappedGroup, seq.Notes[each-1][0])
			}
		}
		groups = append(groups, mappedGroup)
	}
	return core.Sequence{Notes: groups}
}

func NewSequenceMapper(s core.Sequenceable, pattern core.Valueable) SequenceMapper {
	return SequenceMapper{Target: s, Pattern: pattern}
}

func (p SequenceMapper) Storex() string {
	if s, ok := p.Target.(core.Storable); ok {
		if ps, ok := p.Pattern.(core.Storable); ok {
			return fmt.Sprintf("sequencemap(%s,%s)", ps.Storex(), s.Storex())
		}
		return fmt.Sprintf("sequencemap(%v,%s)", p.Pattern, s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (p SequenceMapper) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	if core.IsIdenticalTo(p.Target, from) {
		return SequenceMapper{Target: to, Pattern: p.Pattern}
	}
	if rep, ok := p.Target.(core.Replaceable); ok {
		return SequenceMapper{Target: rep.Replaced(from, to), Pattern: p.Pattern}
	}
	return p
}
