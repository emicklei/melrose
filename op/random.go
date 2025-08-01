package op

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/emicklei/melrose/core"
)

type RandomInteger struct {
	From core.HasValue
	To   core.HasValue
	rnd  *rand.Rand
	last int
}

func NewRandomInteger(from, to core.HasValue) *RandomInteger {
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
	return fmt.Sprintf("random(%s,%s)", core.Storex(r.From), core.Storex(r.To))
}

// Value is part of HasValue
func (r *RandomInteger) Value() any {
	return r.last
}

// Next is part of Nextable
func (r *RandomInteger) Next() any {
	f := core.Int(r.From)
	t := core.Int(r.To)
	if t < f {
		r.last = f
		return f
	}
	r.last = f + r.rnd.Intn(t-f+1)
	return r.last
}

// TODO  Replaceable
