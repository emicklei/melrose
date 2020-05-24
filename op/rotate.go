package op

import (
	"fmt"

	. "github.com/emicklei/melrose"
)

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
