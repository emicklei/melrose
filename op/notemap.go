package op

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

const (
	formatDotAndBangs = iota
	formatNumbers
)

type NoteMap struct {
	Target  core.HasValue
	indices core.HasValue
}

// NewNoteMap returns a NoteMap that creates a sequence from occurrences of a note.
// The format of indices can be one of:
// 1 2 4 ; each number is an index in the sequence where the note is present; rest notes are placed in the gaps.
// ! . ! ; each dot is a rest, each exclamation mark is a presence of a note.
func NewNoteMap(indices core.HasValue, note core.HasValue) NoteMap {
	return NoteMap{
		Target:  note,
		indices: indices,
	}
}

func convertDotsAndBangs(format string) string {
	var b bytes.Buffer
	for i, each := range []rune(format) {
		if each == '.' {
			fmt.Fprintf(&b, "  ")
		} else {
			fmt.Fprintf(&b, "%d ", i+1)
		}
	}
	return b.String()
}

func (n NoteMap) formattedIndices(indices []int, format int, max int) string {
	var b bytes.Buffer
	if format == formatDotAndBangs {
		for i := 1; i <= max; i++ {
			found := false
			for _, each := range indices {
				if each == i {
					found = true
					break
				}
			}
			if found {
				fmt.Fprintf(&b, "!")
			} else {
				fmt.Fprintf(&b, ".")
			}
		}
	} else {
		for i, each := range indices {
			if i > 0 {
				fmt.Fprintf(&b, " ")
			}
			fmt.Fprintf(&b, "%d", each)
		}
	}
	return b.String()
}

func (n NoteMap) Storex() string {
	indexs, format, max := n.getIndices()
	return fmt.Sprintf("notemap('%s',%s)", n.formattedIndices(indexs, format, max), core.Storex(n.Target))
}

func (n NoteMap) getIndices() ([]int, int, int) { // indices, format used, and max index
	indicesStr, ok := core.ValueOf(n.indices).(string)
	if !ok {
		return []int{}, formatNumbers, 0
	}
	idx := []int{}
	// check for dots and bangs first
	var parsed [][]int
	format := formatNumbers
	var maxIndex int
	if strings.ContainsAny(indicesStr, "!.") {
		parsed = parseIndices(convertDotsAndBangs(indicesStr))
		format = formatDotAndBangs
		maxIndex = len(indicesStr)
	} else if strings.ContainsAny(indicesStr, "1234567890 ") { // space is allowed
		parsed = parseIndices(indicesStr)
	} else {
		return []int{}, formatNumbers, 0
	}
	for _, each := range parsed {
		idx = append(idx, each[0])
	}
	max := sliceMax(idx)
	if max > maxIndex {
		maxIndex = max
	}
	return idx, format, maxIndex
}

// Inspect implements Inspectable
func (n NoteMap) Inspect(i core.Inspection) {
	indexs, format, max := n.getIndices()
	if format == formatDotAndBangs {
		i.Properties["nrs"] = n.formattedIndices(indexs, formatNumbers, max)
	} else {
		i.Properties["dots"] = n.formattedIndices(indexs, formatDotAndBangs, max)
	}
	n.S().Inspect(i)
}

func sliceMax(indices []int) int {
	max := 0
	for _, each := range indices {
		if each > max {
			max = each
		}
	}
	return max
}

func (n NoteMap) S() core.Sequence {
	var note core.Note
	notelike, ok := n.Target.Value().(core.NoteConvertable)
	if !ok {
		// try create sequence first
		seq, ok := n.Target.Value().(core.Sequenceable)
		if !ok {
			notify.Console.Errorf("cannot map %v (%T)", n.Target.Value(), n.Target.Value())
			return core.EmptySequence
		}
		// then take the first note
		notes := seq.S()
		if len(notes.Notes) == 0 || len(notes.Notes[0]) == 0 {
			return core.EmptySequence
		}
		note = notes.Notes[0][0]
	} else {
		var err error
		note, err = notelike.ToNote()
		if err != nil {
			notify.Panic(err)
			return core.EmptySequence
		}
	}
	indexs, _, max := n.getIndices()
	notes := make([]core.Note, max)
	for i := range notes {
		notes[i] = note.ToRest()
	}
	for _, each := range indexs {
		notes[each-1] = note
	}
	return core.BuildSequence(notes)
}

// Replaced is part of Replaceable
func (n NoteMap) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(n, from) {
		return to
	}
	notelike, ok := n.Target.Value().(core.NoteConvertable)
	if !ok {
		return n
	}
	note, err := notelike.ToNote()
	if err != nil {
		return n
	}
	return NoteMap{
		Target:  core.On(note.Replaced(from, to)),
		indices: n.indices}
}
