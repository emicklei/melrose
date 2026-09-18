package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type VelocityMap struct {
	Target          []core.Sequenceable
	IndexVelocities []int2int // one-based
}

func NewVelocityMap(target []core.Sequenceable, indices string) VelocityMap {
	return VelocityMap{
		Target:          target,
		IndexVelocities: parseIndexOffsets(indices),
	}
}

func (v VelocityMap) S() core.Sequence {
	if len(v.Target) == 0 {
		return core.EmptySequence
	}
	return core.Sequence{Notes: v.Notes()}
}

func (v VelocityMap) Notes() [][]core.Note {
	if len(v.Target) == 0 {
		return nil
	}
	source := v.Target[0].S().Notes
	target := [][]core.Note{}
	for _, entry := range v.IndexVelocities {
		if entry.from <= 0 || entry.from > len(source) {
			continue
		}
		eachGroup := source[entry.from-1]
		if entry.to == 0 {
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
	if len(v.Target) == 0 {
		return ""
	}
	s, ok := v.Target[0].(core.Storable)
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
	return VelocityMap{Target: replacedAll(v.Target, from, to), IndexVelocities: v.IndexVelocities}
}
