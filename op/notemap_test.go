package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestNewNoteMapper(t *testing.T) {
	m, _ := NewNoteMap("1 2 4", core.On(core.MustParseNote("8c")))
	if got, want := storex(m.S()), "sequence('⅛C ⅛C ⅛= ⅛C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := storex(m), "notemap('1 2 4',note('⅛C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewNoteMapper_Dots(t *testing.T) {
	m, _ := NewNoteMap("!.!.", core.On(core.MustParseNote("c")))
	if got, want := storex(m.S()), "sequence('C = C =')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := storex(m), "notemap('!.!.',note('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNewNoteMapper_MIDI(t *testing.T) {
	mid := core.NewMIDI(core.On(0.5), core.On(60), core.On(60))
	m, _ := NewNoteMap("!.!.", core.On(mid))
	t.Log(m.S())
}
