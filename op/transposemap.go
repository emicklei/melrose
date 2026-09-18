package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type TransposeMap struct {
	indices core.HasValue
	Target  []core.Sequenceable
}

func NewTransposeMap(target []core.Sequenceable, indices core.HasValue) TransposeMap {
	return TransposeMap{
		Target:  target,
		indices: indices,
	}
}

func (p TransposeMap) S() core.Sequence {
	if len(p.Target) == 0 {
		return core.EmptySequence
	}
	return core.Sequence{Notes: p.Notes()}
}

func (p TransposeMap) Notes() [][]core.Note {
	if len(p.Target) == 0 {
		return nil
	}
	source := p.Target[0].S().Notes
	target := [][]core.Note{}
	indicesString, ok := core.ValueOf(p.indices).(string)
	if !ok {
		return nil
	}
	indexOffsets := parseIndexOffsets(indicesString)
	for _, entry := range indexOffsets {
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
			newGroup = append(newGroup, eachNote.Pitched(entry.to))
		}
		target = append(target, newGroup)
	}
	return target
}

// Storex is part of Storable
func (p TransposeMap) Storex() string {
	if len(p.Target) == 0 {
		return ""
	}
	var b bytes.Buffer
	indicesString, ok := core.ValueOf(p.indices).(string)
	if !ok {
		return "?"
	}
	indexOffsets := parseIndexOffsets(indicesString)
	fmt.Fprintf(&b, "transposemap('")
	for i, each := range indexOffsets {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, "%d:%d", each.from, each.to)
	}
	fmt.Fprintf(&b, "',%s", core.Storex(p.Target[0]))
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (p TransposeMap) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(p, from) {
		return to
	}
	return TransposeMap{Target: replacedAll(p.Target, from, to), indices: p.indices}
}
