package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type PitchMap struct {
	IndexOffsets []int2int // one-based
	Target       core.Sequenceable
}

func NewPitchMap(target core.Sequenceable, indices string) PitchMap {
	return PitchMap{
		Target:       target,
		IndexOffsets: parseIndexOffsets(indices),
	}
}

func (p PitchMap) S() core.Sequence {
	return core.Sequence{Notes: p.Notes()}
}

func (p PitchMap) Notes() [][]core.Note {
	source := p.Target.S().Notes
	target := [][]core.Note{}
	for _, entry := range p.IndexOffsets {
		if entry.from <= 0 || entry.from > len(source) {
			// invalid offset, skip
			continue
		}
		eachGroup := source[entry.from-1] // from is one-based
		if entry.to == 0 {
			// no offset, use as is
			target = append(target, eachGroup)
			continue
		}
		newGroup := []core.Note{}
		for _, eachNote := range eachGroup {
			newGroup = append(newGroup, eachNote.Pitched(entry.to))
		}
		target = append(target, newGroup)
	}
	return target
}

// Storex is part of Storable
func (p PitchMap) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "pitchmap('")
	for i, each := range p.IndexOffsets {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, "%d:%d", each.from, each.to)
	}
	fmt.Fprintf(&b, "',%s", core.Storex(p.Target))
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (p PitchMap) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	if core.IsIdenticalTo(p.Target, from) {
		return PitchMap{Target: to, IndexOffsets: p.IndexOffsets}
	}
	if rep, ok := p.Target.(core.Replaceable); ok {
		return PitchMap{Target: rep.Replaced(from, to), IndexOffsets: p.IndexOffsets}
	}
	return p
}
