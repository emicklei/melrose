package op

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/emicklei/melrose/core"
)

/**

prob(0.8,note('c')) =  a sequence with 80% chance of playing C

**/
type Probability struct {
	chance core.HasValue
	seed   *rand.Rand
	target core.HasValue
}

func NewProbability(chance, target core.HasValue) *Probability {
	return &Probability{
		chance: chance,
		seed:   rand.New(rand.NewSource(time.Now().Unix())),
		target: target,
	}
}

func (p *Probability) ToNote() (core.Note, error) {
	v := p.target.Value()
	nc, ok := v.(core.NoteConvertable)
	if !ok {
		return core.Rest4, fmt.Errorf("expect a Note but got %T", v)
	}
	note, err := nc.ToNote()
	if err != nil {
		return core.Rest4, err
	}
	if p.hit() {
		return note, nil
	}
	return note.ToRest(), nil
}

func (p *Probability) S() core.Sequence {
	seq := core.ToSequenceable(p.target).S()
	if p.hit() {
		return seq
	}
	return seq.ToRest()
}

func (p *Probability) hit() bool {
	f := core.Float(p.chance)
	if f > 1 {
		f = f / 100.0
	}
	a := p.seed.Float32()
	return a <= f
}
