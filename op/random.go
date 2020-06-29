package op

import (
	"fmt"
	"github.com/emicklei/melrose/core"
	"math/rand"
	"time"
)

type RandomInteger struct {
	From core.Valueable
	To   core.Valueable
	rnd  *rand.Rand
	last int
}

func NewRandomInteger(from, to core.Valueable) *RandomInteger {
	rnd := &RandomInteger{
		From: from,
		To:   to,
		rnd:  rand.New(rand.NewSource(time.Now().Unix())),
	}
	rnd.Next()
	return rnd
}

// Storex is part of Storable
func (r RandomInteger) Storex() string {
	return fmt.Sprintf("random(%v,%v)", r.From, r.To)
}

// Value is part of Valueable
func (r *RandomInteger) Value() interface{} {
	return r.last
}

// Next is part of Nextable
func (r *RandomInteger) Next() interface{} {
	f := core.Int(r.From)
	t := core.Int(r.To)
	if t < f {
		r.last = f
		return f
	}
	r.last = f + r.rnd.Intn(t-f+1)
	return r.last
}
