package op

import (
	"math/rand"
	"time"

	"github.com/emicklei/melrose/core"
)

/**

prob(0.8,note('c')) =  a sequence with 80% chance of playing C

**/
type Probability struct {
	chance core.Valueable
	seed   *rand.Rand
	target core.Valueable
}

func NewProbability(chance, target core.Valueable) *Probability {
	return &Probability{
		chance: chance,
		seed:   rand.New(rand.NewSource(time.Now().Unix())),
		target: target,
	}
}

func (p *Probability) ToNote() (core.Note, error) {
	if p.hit() {
		v := p.target.Value()
		if n, ok := v.(core.NoteConvertable); ok {
			return n.ToNote()
		}
	}
	return core.Rest4, nil
}

func (p *Probability) S() core.Sequence {
	if p.hit() {
		return core.ToSequenceable(p.target).S()
	}
	return core.EmptySequence
}

func (p *Probability) hit() bool {
	f := core.Float(p.chance)
	if f > 1 {
		f = f / 100.0
	}
	a := p.seed.Float32()
	return a <= f
}
