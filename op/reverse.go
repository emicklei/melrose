package op

import (
	"fmt"
	"github.com/emicklei/melrose/core"
)

type Reverse struct {
	Target core.Sequenceable
}

func (r Reverse) S() core.Sequence {
	return r.Target.S().Reversed()
}

func (r Reverse) Storex() string {
	if s, ok := r.Target.(core.Storable); ok {
		return fmt.Sprintf("reverse(%s)", s.Storex())
	}
	return ""
}
