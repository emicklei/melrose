package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestJoinMapper_S(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	m := NewJoinMap(core.On(j), core.On("(1 2) 1"))
	if got, want := m.Storex(), "joinmap('(1 2 ) 1 ',join(note('C'),note('D')))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoinMapperWithIterator(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	i := &core.Iterator{Target: []any{"1 2"}}
	m := NewJoinMap(core.On(j), i)
	if got, want := m.Storex(), "joinmap('1 2 ',join(note('C'),note('D')))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
