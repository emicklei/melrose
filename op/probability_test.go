package op

import (
	"math/rand"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestProbability_hit(t *testing.T) {
	p := Probability{chance: core.On(0.8), seed: rand.New(rand.NewSource(0))}
	trues := 0
	for i := 0; i < 100; i++ {
		if p.hit() {
			trues++
		}
	}
	t.Log(trues, "out of 100")
	p.chance = core.On(80)
	t.Log(p.hit())
}

func TestProbability_hit_halfrest(t *testing.T) {
	p := Probability{chance: core.On(0.0), seed: rand.New(rand.NewSource(0)), target: core.On(core.N("2c"))}
	n, err := p.ToNote()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := n.DurationFactor(), float32(0.5); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestProbability_hit_quarterrest(t *testing.T) {
	p := Probability{chance: core.On(0.0), seed: rand.New(rand.NewSource(0)),
		target: core.On(core.MustParseSequence("(C e f)"))}
	s := p.S()
	if got, want := s.Storex(), "sequence('(= = =)')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
