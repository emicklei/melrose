package op

import (
	"github.com/emicklei/melrose/core"
	"testing"
)

func TestJoinMapper_S(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	m := NewJoinMapper(core.On(j), "(1 2) 1")
	if got, want := m.S().Storex(), "sequence('(C D) C')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
