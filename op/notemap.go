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
	Target        core.Valueable
	Indices       []int
	indicesFormat int
	maxIndex      int
}

// NewNoteMap returns a NoteMap that creates a sequence from occurrences of a note.
// The format of indices can be one of:
// 1 2 4 ; each number is an index in the sequence where the note is present; rest notes are placed in the gaps.
// ! . ! ; each dot is a rest, each exclamation mark is a presence of a note.
func NewNoteMap(indices string, note core.Valueable) (NoteMap, error) {
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
	if n.indicesFormat == formatDotAndBangs {
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
		return fmt.Sprintf("notemap('%s',%s)", n.formattedIndices(), st.Storex())
	}
	return ""
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
	notelike, ok := n.Target.Value().(core.NoteConvertable)
	if !ok {
		notify.Panic(fmt.Errorf("cannot map"))
		return core.EmptySequence
	}
	notes := make([]core.Note, n.maxIndex)
	note := notelike.ToNote()
	for i := range notes {
		notes[i] = note.ToRest()
	}
	for _, each := range n.Indices {
		notes[each-1] = note
	}
	return core.BuildSequence(notes)
}
