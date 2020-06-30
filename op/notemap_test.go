package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestNewNoteMapper(t *testing.T) {
	m, _ := NewNoteMap("1 2 4", core.On(core.MustParseNote("c")))
	if got, want := storex(m.S()), "sequence('C C = C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewNoteMapper_Dots(t *testing.T) {
	m, _ := NewNoteMap("!.!.", core.On(core.MustParseNote("c")))
	if got, want := storex(m.S()), "sequence('C = C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
