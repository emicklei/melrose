package melrose

import "fmt"

type Join struct {
	Left  Sequenceable
	Right Sequenceable
}

func (j Join) String() string {
	return fmt.Sprintf("(%v).Join(%v)", j.Left, j.Right)
}

func (j Join) Storex() string {
	return fmt.Sprintf("join(%s,%s)", j.Left.Storex(), j.Right.Storex())
}

func (j Join) S() Sequence {
	return j.Left.S().SequenceJoin(j.Right.S())
}
