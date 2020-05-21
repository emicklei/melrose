package melrose

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emicklei/melrose/notify"
)

// TODO  move all operators into package "op"

type Pitch struct {
	Target    Sequenceable
	Semitones Valueable
}

func (p Pitch) S() Sequence {
	return p.Target.S().Pitched(Int(p.Semitones))
}

func (p Pitch) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("pitch(%v,%s)", p.Semitones, s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (p Pitch) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(p, from) {
		return to
	}
	if IsIdenticalTo(p.Target, from) {
		return Pitch{Target: to, Semitones: p.Semitones}
	}
	// https://play.golang.org/p/qHbbK_sTo84
	if r, ok := p.Target.(Replaceable); ok {
		return r.Replaced(from, to)
	}
	return p
}

type Reverse struct {
	Target Sequenceable
}

func (r Reverse) S() Sequence {
	return r.Target.S().Reversed()
}

func (r Reverse) Storex() string {
	if s, ok := r.Target.(Storable); ok {
		return fmt.Sprintf("reverse(%s)", s.Storex())
	}
	return ""
}

type Rotate struct {
	Target Sequenceable
	Times  int
}

func (r Rotate) S() Sequence {
	return r.Target.S().RotatedBy(r.Times)
}

func (r Rotate) Storex() string {
	if s, ok := r.Target.(Storable); ok {
		return fmt.Sprintf("rotate(%d,%s)", r.Times, s.Storex())
	}
	return ""
}

type Serial struct {
	Target []Sequenceable
}

func (a Serial) S() Sequence {
	n := []Note{}
	for _, each := range a.Target {
		each.S().NotesDo(func(each Note) {
			n = append(n, each)
		})
	}
	return BuildSequence(n)
}

func (a Serial) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "serial(")
	for i, each := range a.Target {
		s, ok := each.(Storable)
		if !ok {
			return ""
		}
		fmt.Fprintf(&b, "%s", s.Storex())
		if i < len(a.Target)-1 {
			io.WriteString(&b, ",")
		}
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

type Undynamic struct {
	Target Sequenceable
}

func (u Undynamic) S() Sequence {
	n := []Note{}
	u.Target.S().NotesDo(func(each Note) {
		each.Velocity = Normal
		n = append(n, each)
	})
	return BuildSequence(n)
}

func (u Undynamic) Storex() string {
	if s, ok := u.Target.(Storable); ok {
		return fmt.Sprintf("undynamic(%s)", s.Storex())
	}
	return ""
}

type Parallel struct {
	Target Sequenceable
}

func (p Parallel) S() Sequence {
	n := []Note{}
	p.Target.S().NotesDo(func(each Note) {
		n = append(n, each)
	})
	return Sequence{Notes: [][]Note{n}}
}

func (p Parallel) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("parallel(%s)", s.Storex())
	}
	return ""

}

type SequenceMapper struct {
	Target  Sequenceable
	Indices [][]int
}

func (p SequenceMapper) S() Sequence {
	seq := p.Target.S()
	groups := [][]Note{}
	for _, indexEntry := range p.Indices {
		mappedGroup := []Note{}
		for j, each := range indexEntry {
			if each < 1 || each > len(seq.Notes) {
				notify.Print(notify.Warningf("index out of sequence range: %d, len=%d", j+1, len(seq.Notes)))
			} else {
				// TODO what if sequence had a multi note group?
				mappedGroup = append(mappedGroup, seq.Notes[each-1][0])
			}
		}
		groups = append(groups, mappedGroup)
	}
	return Sequence{Notes: groups}
}

func NewSequenceMapper(s Sequenceable, indices string) SequenceMapper {
	return SequenceMapper{Target: s, Indices: parseIndices(indices)}
}

func (p SequenceMapper) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("sequencemap('%s',%s)", formatIndices(p.Indices), s.Storex())
	}
	return ""
}

// Replaced is part of Replaceable
func (p SequenceMapper) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(p, from) {
		return to
	}
	if IsIdenticalTo(p.Target, from) {
		return SequenceMapper{Target: to, Indices: p.Indices}
	}
	if rep, ok := p.Target.(Replaceable); ok {
		return SequenceMapper{Target: rep.Replaced(from, to), Indices: p.Indices}
	}
	return p
}
