package op

import "github.com/emicklei/melrose/core"

type Stretch struct {
	target []core.Sequenceable
	factor core.Valueable
}

func NewStretch(factor core.Valueable, target []core.Sequenceable) Stretch {
	return Stretch{
		target: target,
		factor: factor,
	}
}

func (s Stretch) S() core.Sequence {
	return Join{Target: s.target}.S().Stretched(core.Float(s.factor))
}
