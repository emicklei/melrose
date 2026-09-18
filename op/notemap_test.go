package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestNewNoteMapper(t *testing.T) {
	m := NewNoteMap(core.On("1 2 4"), core.On(core.MustParseNote("8c")))
	if got, want := storex(m.S()), "sequence('8C 8C 8= 8C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := storex(m), "notemap('1 2 4',note('8C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNoteMap_Storex(t *testing.T) {
	m := NewNoteMap(core.On("1 2 4"), core.On(core.MustParseNote("8c")))
	if got, want := m.Storex(), "notemap('1 2 4',note('8C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNoteMap_Inspect(t *testing.T) {
	m := NewNoteMap(core.On("1 2 4"), core.On(core.MustParseNote("8c")))
	i := core.NewInspect(testContext(), "test", nil)
	i.Properties = map[string]any{}
	m.Inspect(i)
	if got, want := i.Properties["dots"], "!!.!"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	m = NewNoteMap(core.On("!.!"), core.On(core.MustParseNote("8c")))
	i = core.NewInspect(testContext(), "test", nil)
	i.Properties = map[string]any{}
	m.Inspect(i)
	if got, want := i.Properties["nrs"], "1 3"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNoteMap_S(t *testing.T) {
	defer func() { recover() }()
	m := NewNoteMap(core.On("1"), core.On(core.S("c d")))
	if got, want := m.S().Storex(), "sequence('C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	m = NewNoteMap(core.On("1"), core.On(core.S("")))
	if got, want := m.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	m = NewNoteMap(core.On("1"), core.On(failingNoteConvertable{}))
	if got, want := m.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	m = NewNoteMap(core.On("1"), core.On(core.On(1)))
	if got, want := m.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNoteMap_Replaced(t *testing.T) {
	n := core.MustParseNote("c")
	m := NewNoteMap(core.On("1"), core.On(n))
	if core.IsIdenticalTo(m, n) {
		t.Error("should not be identical")
	}
	r := m.Replaced(n, core.S("d"))
	if got, want := r.S().Storex(), "sequence('D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r = m.Replaced(m, core.S("d"))
	if got, want := r.S().Storex(), "sequence('D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	m = NewNoteMap(core.On("1"), core.On(core.On(1)))
	r = m.Replaced(n, core.S("d"))
	if got, want := r.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	// TODO
	// m = NewNoteMap("1", core.On(failingNoteConvertable{}))
	// r = m.Replaced(n, core.S("d"))
	// if got, want := r.S().Storex(), "sequence('')"; got != want {
	// 	t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	// }
}

func TestNewNoteMap_Error(t *testing.T) {
	NewNoteMap(core.On("a"), core.On(core.N("c")))
}

func TestNewNoteMapper_Dots(t *testing.T) {
	m := NewNoteMap(core.On("!.!."), core.On(core.MustParseNote("c")))
	if got, want := storex(m.S()), "sequence('C = C =')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := storex(m), "notemap('!.!.',note('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewNoteMapper_MIDI(t *testing.T) {
	mid := core.NewMIDI(core.On(0.5), core.On(60), core.On(60))
	m := NewNoteMap(core.On("!.!."), core.On(mid))
	t.Log(m.S())
}
