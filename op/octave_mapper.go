package op

import (
	"bytes"
	"fmt"
	"github.com/emicklei/melrose/core"
	"strconv"
	"strings"
)

type OctaveMapper struct {
	Target       core.Sequenceable
	IndexOffsets []int2int // one-based
}

type int2int struct {
	from int
	to   int
}

func NewOctaveMapper(target core.Sequenceable, indices string) OctaveMapper {
	return OctaveMapper{
		Target:       target,
		IndexOffsets: parseIndexOffsets(indices),
	}
}

func (o OctaveMapper) S() core.Sequence {
	return core.Sequence{Notes: o.Notes()}
}

func (o OctaveMapper) Notes() [][]core.Note {
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

// 1:-1,3:-1,1:0,2:0,3:0,1:1,2:1
func parseIndexOffsets(s string) (m []int2int) {
	entries := strings.Split(s, ",")
	for _, each := range entries {
		kv := strings.Split(each, ":")
		ik, err := strconv.Atoi(kv[0])
		if err != nil {
			continue
		}
		iv, err := strconv.Atoi(kv[1])
		if err != nil {
			continue
		}
		m = append(m, int2int{from: ik, to: iv})
	}
	return
}

func (o OctaveMapper) Storex() string {
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
func (o OctaveMapper) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(o, from) {
		return to
	}
	if core.IsIdenticalTo(o.Target, from) {
		return OctaveMapper{Target: to, IndexOffsets: o.IndexOffsets}
	}
	if rep, ok := o.Target.(core.Replaceable); ok {
		return OctaveMapper{Target: rep.Replaced(from, to), IndexOffsets: o.IndexOffsets}
	}
	return o
}
