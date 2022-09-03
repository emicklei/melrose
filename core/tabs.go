package core

import (
	"strconv"
	"strings"
	"time"
)

// tabs('e3 a2 a5 d5 a5 a2 e3')
type BassTablature struct {
	TabNotes []TabNote
}

func (t BassTablature) Storex() string {
	b := new(strings.Builder)
	b.WriteString("tabs('")
	for i, each := range t.TabNotes {
		if i > 0 {
			b.WriteRune(' ')
		}
		each.WriteWith(b)
	}
	b.WriteString("')")
	return b.String()
}

func (t BassTablature) S() Sequence {
	notes := []Note{}
	for _, each := range t.TabNotes {
		notes = append(notes, each.ToNote())
	}
	return BuildSequence(notes)
}

// https://www.studybass.com/lessons/reading-music/how-to-read-bass-tab/
type TabNote struct {
	Name     string // E,A,D,G
	Fret     int    // [0..24]
	Dotted   bool   // if true then fraction is increased by half
	Velocity int    // 1..127

	fraction float32       // {0.0625,0.125,0.25,0.5,1}
	duration time.Duration // if set then this overrides Dotted and fraction
}

func (t TabNote) ToNote() Note {
	n := MakeNote(t.Name, 2, t.fraction, 0, t.Dotted, t.Velocity)
	if t.Name == "D" || t.Name == "G" {
		n = n.Octaved(1)
	}
	return n.Pitched(t.Fret)
}

func (t TabNote) IsRest() bool { return t.Name == "=" }

func (t TabNote) String() string {
	b := new(strings.Builder)
	t.WriteWith(b)
	return b.String()
}

func (t TabNote) WriteWith(b *strings.Builder) {
	if t.fraction != 0.25 {
		b.WriteString(FractionToString(t.fraction))
	}
	if t.Dotted {
		b.WriteRune('.')
	}
	b.WriteString(t.Name)
	if !t.IsRest() {
		b.WriteString(strconv.Itoa(t.Fret))
	}
	if t.Velocity != Normal {
		b.WriteString(VelocityToDynamic(t.Velocity))
	}
}

func ParseTabNote(input string) (TabNote, error) {
	return newFormatParser(input).parseTabNote()
}

func ParseBassTablature(s string) (BassTablature, error) {
	entries := strings.Split(s, " ")
	notes := []TabNote{}
	for _, each := range entries {
		n, err := ParseTabNote(each)
		if err != nil {
			return BassTablature{}, err
		}
		notes = append(notes, n)
	}
	return BassTablature{notes}, nil
}
