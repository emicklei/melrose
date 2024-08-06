package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type ScaleTranspose struct {
	Scale  core.HasValue
	Target core.Sequenceable
	Tones  core.HasValue
}

func (p ScaleTranspose) S() core.Sequence {
	return core.EmptySequence
}

func (s ScaleTranspose) Storex() string {
	return fmt.Sprintf("scale_transpose(%s,%s,%s)", core.Storex(s.Scale), core.Storex(s.Tones), core.Storex(s.Target))
}

// Replaced is part of Replaceable
func (s ScaleTranspose) Replaced(from, to core.Sequenceable) core.Sequenceable {
	return s
}
