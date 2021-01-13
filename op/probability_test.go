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
