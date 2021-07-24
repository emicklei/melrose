package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Rotate struct {
	Target core.Sequenceable
	Times  core.HasValue
}

func (r Rotate) S() core.Sequence {
	return r.Target.S().RotatedBy(core.Int(r.Times))
}

func (r Rotate) Storex() string {
	return fmt.Sprintf("rotate(%s,%s)", core.Storex(r.Times), core.Storex(r.Target))
}
