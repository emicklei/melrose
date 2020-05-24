package op

import (
	"fmt"

	. "github.com/emicklei/melrose"
)

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
