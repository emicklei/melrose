package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Merge struct {
	Target []core.Sequenceable
}

func (m Merge) S() core.Sequence {
	seqs := []core.Sequence{}
	length := 0
	for _, each := range m.Target {
		seq := each.S()
		seqs = append(seqs, seq)
		if seq.Length() > length {
			length = seq.Length()
		}
	}
	merged := [][]core.Note{}
	for i := 0; i < length; i++ {
		group := []core.Note{}
		groupRest := core.Rest4
		for _, each := range seqs {
			if i < each.Length() {
				for _, other := range each.At(i) {
					if !other.IsRest() {
						group = append(group, other)
					} else {
						groupRest = other
					}
				}
			}
		}
		if len(group) == 0 { // only rest notes
			group = append(group, groupRest)
		}
		merged = append(merged, group)
	}
	return core.Sequence{Notes: merged}
}

func (m Merge) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "merge(")
	appendStorexList(&b, true, m.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}
