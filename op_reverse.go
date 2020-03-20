package melrose

import "fmt"

type Reverse struct {
	Target Sequenceable
}

func (r Reverse) S() Sequence {
	return r.Target.S().Reversed()
}

func (r Reverse) String() string {
	return fmt.Sprintf("(%v).Reverse()", r.Target)
}
