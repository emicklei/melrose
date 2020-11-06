package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type OctaveMap struct {
	Target       core.Sequenceable
	IndexOffsets []int2int // one-based
}

type int2int struct {
	from int
	to   int
}

func NewOctaveMap(target core.Sequenceable, indices string) OctaveMap {
	return OctaveMap{
		Target:       target,
		IndexOffsets: parseIndexOffsets(indices),
	}
}

func (o OctaveMap) S() core.Sequence {
	return core.Sequence{Notes: o.Notes()}
}

func (o OctaveMap) Notes() [][]core.Note {
	source := o.Target.S().Notes
	target := [][]core.Note{}
	for _, entry := range o.IndexOffsets {
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
			newGroup = append(newGroup, eachNote.Octaved(entry.to))
		}
		target = append(target, newGroup)
	}
	return target
}

func (o OctaveMap) Storex() string {
	s, ok := o.Target.(core.Storable)
	if !ok {
		return ""
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "octavemap('")
	for i, each := range o.IndexOffsets {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, "%d:%d", each.from, each.to)
	}
	fmt.Fprintf(&b, "',%s", s.Storex())
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (o OctaveMap) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(o, from) {
		return to
	}
	if core.IsIdenticalTo(o.Target, from) {
		return OctaveMap{Target: to, IndexOffsets: o.IndexOffsets}
	}
	if rep, ok := o.Target.(core.Replaceable); ok {
		return OctaveMap{Target: rep.Replaced(from, to), IndexOffsets: o.IndexOffsets}
	}
	return o
}
