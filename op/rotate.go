package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Rotate struct {
	Target []core.Sequenceable
	Times  core.HasValue
}

func (r Rotate) S() core.Sequence {
	if len(r.Target) == 0 {
		return core.EmptySequence
	}
	return r.Target[0].S().RotatedBy(core.Int(r.Times))
}

func (r Rotate) Storex() string {
	if len(r.Target) == 0 {
		return fmt.Sprintf("rotate(%s,nil)", core.Storex(r.Times))
	}
	return fmt.Sprintf("rotate(%s,%s)", core.Storex(r.Times), core.Storex(r.Target[0]))
}
