package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type VelocityMap struct {
	Target          core.Sequenceable
	IndexVelocities []int2int // one-based
}

func NewVelocityMap(target core.Sequenceable, indices string) VelocityMap {
	return VelocityMap{
		Target:          target,
		IndexVelocities: parseIndexOffsets(indices),
	}
}

func (v VelocityMap) S() core.Sequence {
	return core.Sequence{Notes: v.Notes()}
}

func (v VelocityMap) Notes() [][]core.Note {
	source := v.Target.S().Notes
	target := [][]core.Note{}
	for _, entry := range v.IndexVelocities {
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
			newGroup = append(newGroup, eachNote.WithVelocity(entry.to))
		}
		target = append(target, newGroup)
	}
	return target
}

func (v VelocityMap) Storex() string {
	s, ok := v.Target.(core.Storable)
	if !ok {
		return ""
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "velocitymap('")
	for i, each := range v.IndexVelocities {
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
func (v VelocityMap) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(v, from) {
		return to
	}
	if core.IsIdenticalTo(v.Target, from) {
		return OctaveMap{Target: to, IndexOffsets: v.IndexVelocities}
	}
	if rep, ok := v.Target.(core.Replaceable); ok {
		return VelocityMap{Target: rep.Replaced(from, to), IndexVelocities: v.IndexVelocities}
	}
	return v
}
