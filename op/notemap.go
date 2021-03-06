package op

import (
	"bytes"
	"errors"
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
	Target        core.HasValue
	Indices       []int
	indicesFormat int
	maxIndex      int
}

// NewNoteMap returns a NoteMap that creates a sequence from occurrences of a note.
// The format of indices can be one of:
// 1 2 4 ; each number is an index in the sequence where the note is present; rest notes are placed in the gaps.
// ! . ! ; each dot is a rest, each exclamation mark is a presence of a note.
func NewNoteMap(indices string, note core.HasValue) (NoteMap, error) {
	idx := []int{}
	// check for dots and bangs first
	var parsed [][]int
	format := formatNumbers
	var maxIndex int
	if strings.ContainsAny(indices, "!.") {
		parsed = parseIndices(convertDotsAndBangs(indices))
		format = formatDotAndBangs
		maxIndex = len(indices)
	} else if strings.ContainsAny(indices, "1234567890 ") { // space is allowed
		parsed = parseIndices(indices)
	} else {
		return NoteMap{}, errors.New("bad syntax NoteMap; must have digits,spaces OR dots and exclamation marks")
	}
	for _, each := range parsed {
		idx = append(idx, each[0])
	}
	max := sliceMax(idx)
	if max > maxIndex {
		maxIndex = max
	}
	return NoteMap{
		Target:  note,
		Indices: idx,
		// internal
		indicesFormat: format,
		maxIndex:      maxIndex,
	}, nil
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

func (n NoteMap) formattedIndices(format int) string {
	var b bytes.Buffer
	if format == formatDotAndBangs {
		for i := 1; i <= n.maxIndex; i++ {
			found := false
			for _, each := range n.Indices {
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
		for i, each := range n.Indices {
			if i > 0 {
				fmt.Fprintf(&b, " ")
			}
			fmt.Fprintf(&b, "%d", each)
		}
	}
	return b.String()
}

func (n NoteMap) Storex() string {
	st, ok := n.Target.(core.Storable)
	if !ok {
		st, ok = n.Target.Value().(core.Storable)
	}
	if ok {
		return fmt.Sprintf("notemap('%s',%s)", n.formattedIndices(n.indicesFormat), st.Storex())
	}
	return ""
}

// Inspect implements Inspectable
func (n NoteMap) Inspect(i core.Inspection) {
	if n.indicesFormat == formatDotAndBangs {
		i.Properties["nrs"] = n.formattedIndices(formatNumbers)
	} else {
		i.Properties["dots"] = n.formattedIndices(formatDotAndBangs)
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
	notes := make([]core.Note, n.maxIndex)
	for i := range notes {
		notes[i] = note.ToRest()
	}
	for _, each := range n.Indices {
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
		Target:        core.On(note.Replaced(from, to)),
		Indices:       n.Indices,
		indicesFormat: n.indicesFormat,
		maxIndex:      n.maxIndex}
}
