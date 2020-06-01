package op

import (
	"fmt"

	"github.com/emicklei/melrose"
	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

type SequenceMapper struct {
	Target  Sequenceable
	Indices [][]int
	Pattern melrose.Valueable
}

func (p SequenceMapper) S() Sequence {
	if p.Pattern == nil {
		return p.Target.S()
	}
	sPattern := melrose.String(p.Pattern)
	if len(sPattern) == 0 {
		return p.Target.S()
	}
	indices := parseIndices(sPattern)
	seq := p.Target.S()
	groups := [][]Note{}
	for _, indexEntry := range indices {
		mappedGroup := []Note{}
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
	return Sequence{Notes: groups}
}

func NewSequenceMapper(s Sequenceable, pattern melrose.Valueable) SequenceMapper {
	return SequenceMapper{Target: s, Pattern: pattern}
}

func (p SequenceMapper) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("sequencemap(%v,%s)", p.Pattern, s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (p SequenceMapper) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(p, from) {
		return to
	}
	if IsIdenticalTo(p.Target, from) {
		return SequenceMapper{Target: to, Pattern: p.Pattern}
	}
	if rep, ok := p.Target.(Replaceable); ok {
		return SequenceMapper{Target: rep.Replaced(from, to), Pattern: p.Pattern}
	}
	return p
}
