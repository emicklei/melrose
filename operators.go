package melrose

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emicklei/melrose/notify"
)

type Pitch struct {
	Target    Sequenceable
	Semitones int
}

func (p Pitch) S() Sequence {
	return p.Target.S().Pitched(p.Semitones)
}

func (p Pitch) Storex() string {
	return fmt.Sprintf("pitch(%d,%s)", p.Semitones, p.Target.Storex())
}

type Join struct {
	List []Sequenceable
}

func (j Join) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "join(")
	for i, each := range j.List {
		fmt.Fprintf(&b, "%s", each.Storex())
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
	return fmt.Sprintf("repeat(%d,%s)", r.Times, r.Target.Storex())
}

type Reverse struct {
	Target Sequenceable
}

func (r Reverse) S() Sequence {
	return r.Target.S().Reversed()
}

func (r Reverse) Storex() string {
	return fmt.Sprintf("reverse(%s)", r.Target.Storex())
}

type Rotate struct {
	Target Sequenceable
	Times  int
}

func (r Rotate) S() Sequence {
	return r.Target.S().RotatedBy(r.Times)
}

func (r Rotate) Storex() string {
	return fmt.Sprintf("rotate(%d,%s)", r.Times, r.Target.Storex())
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
	return fmt.Sprintf("serial(%s)", a.Target.Storex())
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
	return fmt.Sprintf("undynamic(%s)", u.Target.Storex())
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
	return fmt.Sprintf("parallel(%s)", p.Target.Storex())
}

type IndexMapper struct {
	Target  Sequenceable
	Indices []int
}

func (p IndexMapper) S() Sequence {
	seq := p.Target.S()
	groups := [][]Note{}
	for j, i := range p.Indices {
		if i < 0 || i > len(seq.Notes) {
			notify.Print(notify.Warningf("index out of sequence range: %d=%d", j, i))
		} else {
			groups = append(groups, seq.Notes[i-1])
		}
	}
	return Sequence{Notes: groups}
}

func (p IndexMapper) Storex() string {
	return fmt.Sprintf("indexmap(%s)", p.Target.Storex())
}
