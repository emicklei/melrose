package op

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/emicklei/melrose/core"
)

const (
	formatDotAndBangs = iota
	formatDigits
)

type NoteMap struct {
	Target        core.Valueable
	Indices       []int
	IndicesFormat int
}

// NewNoteMap returns a NoteMap that creates a sequence from occurrences of a note.
// The format of indices can be one of:
// 1 2 4 ; each number is an index in the sequence where the note is present; rest notes are placed in the gaps.
// ! . ! ; each dot is a rest, each exclamation mark is a presence of a note.
func NewNoteMap(indices string, note core.Valueable) (NoteMap, error) {
	idx := []int{}
	// check for dots and bangs first
	var parsed [][]int
	format := formatDigits
	if strings.ContainsAny(indices, "!.") {
		parsed = parseIndices(convertDotsAndBangs(indices))
		format = formatDotAndBangs
	} else if strings.ContainsAny(indices, "1234567890 ") {
		parsed = parseIndices(indices)
	} else {
		return NoteMap{}, errors.New("bad syntax NoteMap; must have digits,spaces OR dots and exclamation marks")
	}
	for _, each := range parsed {
		idx = append(idx, each[0])
	}
	return NoteMap{
		Target:        note,
		Indices:       idx,
		IndicesFormat: format,
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

func (n NoteMap) formattedIndices() string {
	var b bytes.Buffer
	if n.IndicesFormat == formatDotAndBangs {
		fmt.Fprintf(&b, "?")
	}
	return b.String()
}

func (n NoteMap) Storex() string {
	if st, ok := n.Target.Value().(core.Storable); ok {
		return fmt.Sprintf("notemap('%s',%s)", n.formattedIndices(), st.Storex())
	}
	return ""
}

func (n NoteMap) minMax() (int, int) {
	max := 0
	min := 10000
	for _, each := range n.Indices {
		if each > max {
			max = each
		} else if each < min {
			min = each
		}
	}
	return min, max
}

func (n NoteMap) S() core.Sequence {
	notelike, ok := n.Target.Value().(core.NoteConvertable)
	if !ok {
		// TODO warning here?
		return core.EmptySequence
	}
	_, max := n.minMax()
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
