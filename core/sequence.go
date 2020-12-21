package core

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"time"
)

var EmptySequence = Sequence{}

type Sequence struct {
	Notes [][]Note
}

// At uses zero-based indexing into the Notes
func (s Sequence) At(i int) []Note {
	if i < 0 || i > len(s.Notes)-1 {
		panic("Sequence index out of bounds:" + strconv.Itoa(i))
	}
	return s.Notes[i]
}

// Replaced is part of Replaceable
func (s Sequence) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, s) {
		return to
	}
	return s
}

// SequenceJoin returns s + t
func (s Sequence) SequenceJoin(t Sequence) Sequence {
	return Sequence{append(s.Notes, t.Notes...)}
}

func (s Sequence) NotesDo(block func(Note)) {
	for _, eachGroup := range s.Notes {
		for _, eachNote := range eachGroup {
			block(eachNote)
		}
	}
}

// BuildSequence creates a Sequence from a slice of Note
func BuildSequence(notes []Note) Sequence {
	groups := [][]Note{}
	for _, each := range notes {
		groups = append(groups, []Note{each})
	}
	return Sequence{Notes: groups}
}

var S = MustParseSequence

func MustParseSequence(input string) Sequence {
	if s, err := ParseSequence(input); err != nil {
		log.Fatal("MustParseSequence failed:", err.Error())
		return s
	} else {
		return s
	}
}

// ParseSequence creates a Sequence by reading the format "Note* [Note Note*]* Note*"
func ParseSequence(input string) (Sequence, error) {
	return newFormatParser(input).parseSequence()
}

func (s Sequence) S() Sequence {
	return s
}

func (s Sequence) DurationFactor() float64 {
	dur := float32(0.0)
	for _, each := range s.Notes {
		if len(each) > 0 {
			lead := each[0]
			dur += lead.DurationFactor()
		}
	}
	return float64(dur)
}

func (s Sequence) Inspect(i Inspection) {
	i.Properties["duration"] = s.Duration(i.Context.Control().BPM())
	i.Properties["note(s)/groups"] = len(s.Notes)
	i.Properties["bars"] = s.Bars(i.Context.Control().BIAB())
}

// Conversion

// Storex returns the command line expression that creates the receiver
func (s Sequence) Storex() string {
	return fmt.Sprintf("sequence('%s')", s.String())
}

func (s Sequence) String() string {
	return s.PrintString(PrintAsSpecified)
}

func (s Sequence) PrintString(sharpOrFlatKey int) string {
	var buf bytes.Buffer
	s.writeNotesOn(&buf, (Note).printOn, sharpOrFlatKey)
	return buf.String()
}

const (
	groupOpen  = "("
	groupClose = ")"
)

func (s Sequence) writeNotesOn(
	buf *bytes.Buffer,
	printer func(n Note, buf *bytes.Buffer, sharpOrFlatKey int),
	sharpOrFlatKey int) {

	for i, each := range s.Notes {
		if i > 0 {
			buf.WriteString(" ")
		}
		if len(each) > 1 {
			buf.WriteString(groupOpen)
		}
		for j, other := range each {
			if j > 0 {
				buf.WriteString(" ")
			}
			printer(other, buf, sharpOrFlatKey)
		}
		if len(each) > 1 {
			buf.WriteString(groupClose)
		}
	}
}

func StringFromNoteGroup(notes []Note) string {
	var buf bytes.Buffer
	buf.WriteString(groupOpen)
	for i, each := range notes {
		if i > 0 {
			buf.WriteString(" ")
		}
		each.printOn(&buf, PrintAsSpecified)
	}
	buf.WriteString(groupClose)
	return buf.String()
}

func (s Sequence) Duration(bpm float64) time.Duration {
	w := WholeNoteDuration(bpm)
	l := time.Duration(0)
	for _, group := range s.Notes {
		if len(group) > 0 {
			l += time.Duration(float32(w) * group[0].DurationFactor())
		}
	}
	return l
}

func (s Sequence) Bars(biab int) float64 {
	return float64(s.DurationFactor()) * 4 / float64(biab) // 4 because signature denominator
}

// W returns the mapping of each note to a delta of semitones compared to middle C4.
// Can be used for the pitch lane of the Korg Wavestate
func (s Sequence) W() string {
	var buf bytes.Buffer
	ref, _ := NewNote("C", 4, 0.25, 0, false, Normal)
	mapit := func(n Note) {
		fmt.Fprintf(&buf, "%s:%d", n.String(), n.MIDI()-ref.MIDI())
	}
	for i, group := range s.Notes {
		if i > 0 {
			buf.WriteString(" ")
		}
		if len(group) > 1 {
			buf.WriteString(groupOpen)
			for i, note := range group {
				if i > 0 {
					buf.WriteString(" ")
				}
				if note.IsHearable() {
					mapit(note)
				}
			}
			buf.WriteString(groupClose)
		} else {
			if len(group) == 1 && group[0].IsHearable() {
				mapit(group[0])
			}
		}
	}
	return buf.String()
}
