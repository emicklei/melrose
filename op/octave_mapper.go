package op

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/emicklei/melrose"
)

type OctaveMapper struct {
	Target       melrose.Sequenceable
	IndexOffsets []int2int // one-based
}

type int2int struct {
	from int
	to   int
}

func NewOctaveMapper(target melrose.Sequenceable, indices string) OctaveMapper {
	return OctaveMapper{
		Target:       target,
		IndexOffsets: parseIndexOffsets(indices),
	}
}

func (o OctaveMapper) S() melrose.Sequence {
	return melrose.Sequence{Notes: o.Notes()}
}

func (o OctaveMapper) Notes() [][]melrose.Note {
	source := o.Target.S().Notes
	target := [][]melrose.Note{}
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
		newGroup := []melrose.Note{}
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
	s, ok := o.Target.(melrose.Storable)
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
	fmt.Fprintf(&b, s.Storex())
	fmt.Fprintf(&b, ")")
	return b.String()
}
