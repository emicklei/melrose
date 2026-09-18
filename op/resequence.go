package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
)

type Resequencer struct {
	Target  []core.Sequenceable
	Indices [][]int
	Pattern core.HasValue
}

func (p Resequencer) S() core.Sequence {
	if len(p.Target) == 0 {
		return core.EmptySequence
	}
	if p.Pattern == nil {
		return p.Target[0].S()
	}
	sPattern := core.String(p.Pattern)
	if len(sPattern) == 0 {
		return p.Target[0].S()
	}
	indices := parseIndices(sPattern)
	seq := p.Target[0].S()
	groups := [][]core.Note{}
	for _, indexEntry := range indices {
		mappedGroup := []core.Note{}
		for j, each := range indexEntry {
			if each < 1 || each > len(seq.Notes) {
				notify.Warnf("index out of sequence range: %d, len=%d", j+1, len(seq.Notes))
			} else {
				mappedGroup = append(mappedGroup, seq.Notes[each-1][0])
			}
		}
		groups = append(groups, mappedGroup)
	}
	return core.Sequence{Notes: groups}
}

func NewResequencer(s []core.Sequenceable, pattern core.HasValue) Resequencer {
	return Resequencer{Target: s, Pattern: pattern}
}

func (p Resequencer) Storex() string {
	if len(p.Target) == 0 {
		return "?"
	}
	if s, ok := p.Target[0].(core.Storable); ok {
		if ps, ok := p.Pattern.(core.Storable); ok {
			return fmt.Sprintf("resequence(%s,%s)", ps.Storex(), s.Storex())
		}
		return "?"
	}
	return "?"
}

// Replaced is part of Replaceable
func (p Resequencer) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	return Resequencer{Target: replacedAll(p.Target, from, to), Pattern: p.Pattern}
}
