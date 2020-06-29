package op

import (
	"fmt"
	"github.com/emicklei/melrose/core"
)

type Rotate struct {
	Target core.Sequenceable
	Times  int
}

func (r Rotate) S() core.Sequence {
	return r.Target.S().RotatedBy(r.Times)
}

func (r Rotate) Storex() string {
	if s, ok := r.Target.(core.Storable); ok {
		return fmt.Sprintf("rotate(%d,%s)", r.Times, s.Storex())
	}
	return ""
}
