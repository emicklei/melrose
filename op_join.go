package melrose

import (
	"bytes"
	"fmt"
	"io"
)

type Join struct {
	List []Sequenceable
}

func (j Join) String() string {
	if len(j.List) == 0 {
		return ""
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "(%v).Join(", j.List[0])
	for i, each := range j.List {
		if i > 0 {
			fmt.Fprintf(&b, "%s", each.Storex())
			if i < len(j.List)-1 {
				io.WriteString(&b, ",")
			}
		}
	}
	fmt.Fprintf(&b, ")")
	return b.String()
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
