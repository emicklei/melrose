package melrose

import (
	"bytes"
	"strings"
	"time"
)

type Sequence struct {
	Notes [][]Note
}

func (s Sequence) Size() int {
	sum := 0
	s.NotesDo(func(n Note) {
		sum++
	})
	return sum
}

// Append creates a new Sequence with Notes at the end.
func (s Sequence) Append(notes ...Note) Sequence {
	list := s.Notes
	for _, each := range notes {
		list = append(list, []Note{each})
	}
	return Sequence{list}
}

func (s Sequence) Join(others ...Joinable) Sequence {
	if len(others) == 0 {
		return s
	}
	if len(others) == 1 {
		return others[0].SequenceJoin(s)
	}
	return others[0].SequenceJoin(s).Join(others[1:]...)
}

// SequenceJoin returns t + s
func (s Sequence) SequenceJoin(t Sequence) Sequence {
	return Sequence{append(t.Notes, s.Notes...)}
}

// NoteJoin returns n + s
func (s Sequence) NoteJoin(n Note) Sequence {
	return Sequence{append([][]Note{[]Note{n}}, s.Notes...)}
}

func (s Sequence) NotesCollect(transform func(Note) Note) Sequence {
	notes := make([][]Note, len(s.Notes))
	for i, eachGroup := range s.Notes {
		group := make([]Note, len(eachGroup))
		for j, eachNote := range eachGroup {
			note := transform(eachNote)
			group[j] = note
		}
		notes[i] = group
	}
	return Sequence{Notes: notes}
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
		panic("MustParseSequence failed:" + err.Error())
	} else {
		return s
	}

}

// ParseSequence creates a Sequence by reading the format "Note* (Note Note*)* Note*"
func ParseSequence(input string) (Sequence, error) {
	m := Sequence{}
	// hack to keep scanning simple, TODO
	splitable := strings.Replace(input, "(", " ( ", -1)
	splitable = strings.Replace(splitable, ")", " ) ", -1)
	parts := strings.Fields(splitable)
	ingroup := false
	var group []Note
	for _, each := range parts {
		if "(" == each {
			ingroup = true
			group = []Note{}
		} else if ")" == each {
			ingroup = false
			m.Notes = append(m.Notes, group)
		} else {
			next, err := ParseNote(each)
			if err != nil {
				return m, err
			}
			if ingroup {
				group = append(group, next)
			} else {
				m.Notes = append(m.Notes, []Note{next})
			}
		}
	}
	return m, nil
}

func (s Sequence) Repeated(howMany int) Sequence {
	r := s
	for i := 0; i < howMany; i++ {
		r = r.Join(s)
	}
	return r
}

func (s Sequence) RotatedBy(direction int, howMany int) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	groups := s.Notes
	for c := 0; c < howMany; c++ {
		if direction == Left {
			first := groups[0]
			groups = append(groups[1:], first)
		} else {
			last := groups[len(groups)-1]
			groups = append([][]Note{last}, groups[:len(groups)-1]...)
		}
	}
	return Sequence{groups}
}

func (s Sequence) Reversed() Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	groups := [][]Note{}
	for c := len(s.Notes) - 1; c != -1; c-- {
		groups = append(groups, s.Notes[c])
	}
	return Sequence{groups}
}

// Conversion

func (s Sequence) String() string {
	return s.PrintString(PrintAsSpecified)
}

func (s Sequence) PrintString(sharpOrFlatKey int) string {
	var buf bytes.Buffer
	s.writeNotesOn(&buf, (Note).PrintOn, sharpOrFlatKey)
	return buf.String()
}

func (s Sequence) EncodeOn(buf *bytes.Buffer, sharpOrFlatKey int) {
	s.writeNotesOn(buf, (Note).EncodeOn, sharpOrFlatKey)
}

func (s Sequence) writeNotesOn(
	buf *bytes.Buffer,
	printer func(n Note, buf *bytes.Buffer, sharpOrFlatKey int),
	sharpOrFlatKey int) {

	for i, each := range s.Notes {
		if i > 0 {
			buf.WriteString(" ")
		}
		if len(each) > 1 {
			buf.WriteString("(")
		}
		for j, other := range each {
			if j > 0 {
				buf.WriteString(" ")
			}
			printer(other, buf, sharpOrFlatKey)
		}
		if len(each) > 1 {
			buf.WriteString(")")
		}
	}
}

func (s Sequence) Play(p Player, t time.Duration) {
	p.PlaySequence(s, t)
}
