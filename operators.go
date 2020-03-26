package melrose

import (
	"bytes"
	"fmt"
	"io"
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
	return fmt.Sprintf("reverse(%s)", r.Target)
}

type Rotate struct {
	Target Sequenceable
	Times  int
}

func (r Rotate) S() Sequence {
	return r.Target.S().RotatedBy(r.Times)
}

func (r Rotate) Storex() string {
	return fmt.Sprintf("rotate(%d,%s)", r.Times, r.Target)
}

type Ungroup struct {
	Target Sequenceable
}

func (a Ungroup) S() Sequence {
	n := []Note{}
	a.Target.S().NotesDo(func(each Note) {
		n = append(n, each)
	})
	return BuildSequence(n)
}

func (a Ungroup) Storex() string {
	return fmt.Sprintf("flat(%s)", a.Target)
}
