package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestJoinMapper_S(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	m := NewJoinMap(core.On(j), "(1 2) 1")
	if got, want := m.Storex(), "joinmap('(1 2 ) 1 ',join(note('C'),note('D')))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
