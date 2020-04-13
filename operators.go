package melrose

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emicklei/melrose/notify"
)

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

type Join struct {
	List []Sequenceable
}

func (j Join) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "join(")
	for i, each := range j.List {
		s, ok := each.(Storable)
		if !ok {
			return ""
		}
		fmt.Fprintf(&b, "%s", s.Storex())
		if i < len(j.List)-1 {
			io.WriteString(&b, ",")
		}
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (j Join) S() Sequence {
	if len(j.List) == 0 {
		return Sequence{}
	}
	head := j.List[0].S()
	for i := 1; i < len(j.List); i++ {
		head = head.SequenceJoin(j.List[i].S())
	}
	return head
}

type Repeat struct {
	Target Sequenceable
	Times  int
}

func (r Repeat) S() Sequence {
	return r.Target.S().Repeated(r.Times)
}

func (r Repeat) Storex() string {
	if s, ok := r.Target.(Storable); ok {
		return fmt.Sprintf("repeat(%d,%s)", r.Times, s.Storex())
	}
	return ""
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
	Target Sequenceable
}

func (a Serial) S() Sequence {
	n := []Note{}
	a.Target.S().NotesDo(func(each Note) {
		n = append(n, each)
	})
	return BuildSequence(n)
}

func (a Serial) Storex() string {
	if s, ok := a.Target.(Storable); ok {
		return fmt.Sprintf("serial(%s)", s.Storex())
	}
	return ""
}

type Undynamic struct {
	Target Sequenceable
}

func (u Undynamic) S() Sequence {
	n := []Note{}
	u.Target.S().NotesDo(func(each Note) {
		each.velocityFactor = 1.0
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

type IndexMapper struct {
	Target  Sequenceable
	Indices [][]int
}

func (p IndexMapper) S() Sequence {
	seq := p.Target.S()
	groups := [][]Note{}
	for _, group := range p.Indices {
		mappedGroup := []Note{}
		for j, each := range group {
			if each < 0 || each > len(seq.Notes) {
				notify.Print(notify.Warningf("index out of sequence range: %d=%d", j, each))
			} else {
				// TODO what if sequence had a multi note group?
				mappedGroup = append(mappedGroup, seq.Notes[each-1][0])
			}
		}
		groups = append(groups, mappedGroup)
	}
	return Sequence{Notes: groups}
}

func NewIndexMapper(s Sequenceable, indices string) IndexMapper {
	return IndexMapper{Target: s, Indices: parseIndices(indices)}
}

func (p IndexMapper) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("indexmap(%s)", s.Storex())
	}
	return ""
}
