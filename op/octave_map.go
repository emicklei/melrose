package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type OctaveMap struct {
	Target  []core.Sequenceable
	indices core.HasValue
}

func NewOctaveMap(target []core.Sequenceable, indices core.HasValue) OctaveMap {
	return OctaveMap{
		Target:  target,
		indices: indices,
	}
}

func (o OctaveMap) S() core.Sequence {
	if len(o.Target) == 0 {
		return core.EmptySequence
	}
	return core.Sequence{Notes: o.Notes()}
}

func (o OctaveMap) Notes() [][]core.Note {
	if len(o.Target) == 0 {
		return nil
	}
	source := o.Target[0].S().Notes
	target := [][]core.Note{}
	indicesString, ok := core.ValueOf(o.indices).(string)
	if !ok {
		return nil
	}
	for _, entry := range parseIndexOffsets(indicesString) {
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
			newGroup = append(newGroup, eachNote.Octaved(entry.to))
		}
		target = append(target, newGroup)
	}
	return target
}

func (o OctaveMap) Storex() string {
	if len(o.Target) == 0 {
		return ""
	}
	s, ok := o.Target[0].(core.Storable)
	if !ok {
		return ""
	}
	var b bytes.Buffer
	indicesString, ok := core.ValueOf(o.indices).(string)
	if !ok {
		return ""
	}
	fmt.Fprintf(&b, "octavemap('")
	for i, each := range parseIndexOffsets(indicesString) {
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
	return OctaveMap{Target: replacedAll(o.Target, from, to), indices: o.indices}
}
