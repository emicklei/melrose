package op

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emicklei/melrose/core"
	"strings"
)

type NoteMap struct {
	Target  core.Valueable
	Indices []int
}

// NewNoteMapper returns a NoteMap that creates a sequence from occurrences of a note.
// The format of indices can be one of:
// 1 2 4 ; each number is an index in the sequence where the note is present; rest notes are placed in the gaps.
// ! . ! ; each dot is a rest, each exclamation mark is a presence of a note.
func NewNoteMapper(indices string, note core.Valueable) (NoteMap, error) {
	idx := []int{}
	// check for dots and bangs first
	var parsed [][]int
	if strings.ContainsAny(indices, "!.") {
		parsed = parseIndices(convertDotsAndBangs(indices))
	} else if strings.ContainsAny(indices, "1234567890 ") {
		parsed = parseIndices(indices)
	} else {
		return NoteMap{}, errors.New("bad syntax NoteMap; must have digits,spaces OR dots and exclamation marks")
	}
	for _, each := range parsed {
		idx = append(idx, each[0])
	}
	return NoteMap{
		Target:  note,
		Indices: idx,
	}, nil
}

func convertDotsAndBangs(format string) string {
	var b bytes.Buffer
	for i, each := range []rune(format) {
		if '.' == each {
			fmt.Fprintf(&b, "  ")
		} else {
			fmt.Fprintf(&b, "%d ", i+1)
		}
	}
	return b.String()
}

func (n NoteMap) S() core.Sequence {
	notelike, ok := n.Target.Value().(core.NoteConvertable)
	if !ok {
		// TODO warning here?
		return core.EmptySequence
	}
	max := 0
	min := 10000
	for _, each := range n.Indices {
		if each > max {
			max = each
		} else if each < min {
			min = each
		}
	}
	notes := make([]core.Note, max)
	for i := range notes {
		notes[i] = core.Rest4
	}
	note := notelike.ToNote()
	for _, each := range n.Indices {
		notes[each-1] = note
	}
	return core.BuildSequence(notes)
}

type NoteMerge struct {
	Target []core.Valueable
	Count  int
}

func NewNoteMerge(count int, maps []core.Valueable) NoteMerge {
	return NoteMerge{
		Count:  count,
		Target: maps,
	}
}

var restGroup = []core.Note{core.Rest4}

func (m NoteMerge) S() core.Sequence {
	groups := [][]core.Note{}
	for g := 1; g <= m.Count; g++ {
		group := []core.Note{}
		for _, eachMapVal := range m.Target {
			eachMap, ok := eachMapVal.Value().(NoteMap)
			if !ok {
				// TODO warning here?
				return core.EmptySequence
			}
			for _, eachIndex := range eachMap.Indices {
				if eachIndex == g {
					notelike, ok := eachMap.Target.Value().(core.NoteConvertable)
					if !ok {
						// TODO warning here?
						return core.EmptySequence
					}
					group = append(group, notelike.ToNote())
					break
				}
			}
		}
		if len(group) == 0 {
			groups = append(groups, restGroup)
		} else {
			groups = append(groups, group)
		}
	}
	return core.Sequence{Notes: groups}
}
