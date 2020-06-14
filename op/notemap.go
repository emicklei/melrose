package op

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/emicklei/melrose"
)

type NoteMap struct {
	Target  melrose.Valueable
	Indices []int
}

// NewNoteMapper returns a NoteMap that creates a sequence from occurrences of a note.
// The format of indices can be one of:
// 1 2 4 ; each number is an index in the sequence where the note is present; rest notes are placed in the gaps.
// ! . ! ; each dot is a rest, each exclamation mark is a presence of a note.
func NewNoteMapper(indices string, note melrose.Valueable) (NoteMap, error) {
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

func (n NoteMap) S() melrose.Sequence {
	notelike, ok := n.Target.Value().(melrose.NoteConvertable)
	if !ok {
		// TODO warning here?
		return melrose.EmptySequence
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
	notes := make([]melrose.Note, max)
	for i := range notes {
		notes[i] = melrose.Rest4
	}
	note := notelike.ToNote()
	for _, each := range n.Indices {
		notes[each-1] = note
	}
	return melrose.BuildSequence(notes)
}

type NoteMerge struct {
	Target []melrose.Valueable
	Count  int
}

func NewNoteMerge(count int, maps []melrose.Valueable) NoteMerge {
	return NoteMerge{
		Count:  count,
		Target: maps,
	}
}

var restGroup = []melrose.Note{melrose.Rest4}

func (m NoteMerge) S() melrose.Sequence {
	groups := [][]melrose.Note{}
	for g := 1; g <= m.Count; g++ {
		group := []melrose.Note{}
		for _, eachMapVal := range m.Target {
			eachMap, ok := eachMapVal.Value().(NoteMap)
			if !ok {
				// TODO warning here?
				return melrose.EmptySequence
			}
			for _, eachIndex := range eachMap.Indices {
				if eachIndex == g {
					notelike, ok := eachMap.Target.Value().(melrose.NoteConvertable)
					if !ok {
						// TODO warning here?
						return melrose.EmptySequence
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
	return melrose.Sequence{Notes: groups}
}
