package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestNewNoteMapper(t *testing.T) {
	m, _ := NewNoteMapper("1 2 4", melrose.On(melrose.MustParseNote("c")))
	if got, want := m.S().String(), "C C = C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewNoteMapper_Dots(t *testing.T) {
	m, _ := NewNoteMapper("!.!.", melrose.On(melrose.MustParseNote("c")))
	if got, want := m.S().String(), "C = C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
