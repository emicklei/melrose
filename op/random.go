package op

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/emicklei/melrose"
)

type RandomInteger struct {
	From melrose.Valueable
	To   melrose.Valueable
	rnd  *rand.Rand
	last int
}

func NewRandomInteger(from, to melrose.Valueable) *RandomInteger {
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
	f := melrose.Int(r.From)
	t := melrose.Int(r.To)
	if t < f {
		r.last = f
		return f
	}
	r.last = f + r.rnd.Intn(t-f+1)
	return r.last
}
