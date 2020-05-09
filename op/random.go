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
}

func NewRandomInteger(from, to melrose.Valueable) *RandomInteger {
	return &RandomInteger{
		From: from,
		To:   to,
		rnd:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (r RandomInteger) Storex() string {
	return fmt.Sprintf("random(%v,%v)", r.From, r.To)
}

func (r *RandomInteger) Value() interface{} {
	f := melrose.Int(r.From)
	t := melrose.Int(r.To)
	if t < f {
		return f
	}
	return f + r.rnd.Intn(t-f+1)
}
