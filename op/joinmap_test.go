package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestJoinMapper_S(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	m := NewJoinMap(core.On(j), core.On("(1 2) 1"))
	if got, want := m.Storex(), "joinmap('(1 2) 1',join(note('C'),note('D')))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoinMapperWithIterator(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	i := &core.Iterator{Target: []any{"1 2"}}
	m := NewJoinMap(core.On(j), i)
	if got, want := m.Storex(), "joinmap('1 2',join(note('C'),note('D')))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoinMap_S(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	m := NewJoinMap(core.On(j), core.On("(1 2) 1"))
	if got, want := m.S().Storex(), "sequence('(C D) C')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	m = NewJoinMap(core.On(core.S("a b")), core.On("1"))
	if got, want := m.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	m = NewJoinMap(core.On(j), core.On("3"))
	if got, want := m.S().Storex(), "sequence('=')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	m = NewJoinMap(core.On(j), core.On("(1 3)"))
	if got, want := m.S().Storex(), "sequence('(C =)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoinMap_Replaced(t *testing.T) {
	j := Join{Target: []core.Sequenceable{core.MustParseNote("c"), core.MustParseNote("d")}}
	m := NewJoinMap(core.On(j), core.On("(1 2) 1"))
	if core.IsIdenticalTo(m, j) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(m.Replaced(j, core.EmptySequence).(JoinMap).target.Value().(core.Sequenceable), core.EmptySequence) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(m.Replaced(m, j), j) {
		t.Error("should be replaced by j")
	}
	m = NewJoinMap(core.On(core.S("a b")), core.On("1"))
	if !core.IsIdenticalTo(m.Replaced(m, j), j) {
		t.Error("should be replaced by j")
	}
}
