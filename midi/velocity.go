package midi

import (
	"math/rand"
	"time"
)

type VelocityModifier interface {
	Offset() int
}

type VelocityOffset struct {
	min int
	max int
	rnd *rand.Rand
}

func newVelocityOffset(min, max int) VelocityOffset {
	return VelocityOffset{
		min: min,
		max: max,
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (v VelocityOffset) Offset() int {
	return v.min + v.rnd.Intn(v.max-v.min)
}
