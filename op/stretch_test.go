package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestStretch_S(t *testing.T) {
	s := core.MustParseSequence("(c d) 2e 8.f")
	st := s.Stretched(2.0)
	if got, want := core.Storex(st), "sequence('(2C 2D) 1E .F')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestStretch_Storex(t *testing.T) {
	s := NewStretch(core.On(2.0), []core.Sequenceable{core.MustParseSequence("C")})
	if got, want := s.Storex(), "stretch(2,sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestStretch_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	s := NewStretch(core.On(2.0), []core.Sequenceable{s1})
	if core.IsIdenticalTo(s, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(s.Replaced(s1, s2).(Stretch).target[0], s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(s.Replaced(s, s2), s2) {
		t.Error("should be replaced by s2")
	}
}
