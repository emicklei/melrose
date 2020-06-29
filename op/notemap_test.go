package op

import (
	"github.com/emicklei/melrose/core"
	"testing"
)

func TestNewNoteMapper(t *testing.T) {
	m, _ := NewNoteMapper("1 2 4", core.On(core.MustParseNote("c")))
	if got, want := m.S().String(), "C C = C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewNoteMapper_Dots(t *testing.T) {
	m, _ := NewNoteMapper("!.!.", core.On(core.MustParseNote("c")))
	if got, want := m.S().String(), "C = C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
