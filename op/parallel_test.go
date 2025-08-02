package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestGroup_S(t *testing.T) {
	s := core.MustParseSequence("C D E")
	g := Group{Target: s}
	if got, want := g.S().Storex(), "sequence('(C D E)')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestGroup_Storex(t *testing.T) {
	s := core.MustParseSequence("C D E")
	g := Group{Target: s}
	if got, want := g.Storex(), "group(sequence('C D E'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestGroup_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	g := Group{Target: s1}
	if core.IsIdenticalTo(g, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(g.Replaced(s1, s2).(Group).Target, s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(g.Replaced(g, s2), s2) {
		t.Error("should be replaced by s2")
	}
	g = Group{Target: g}
	if !core.IsIdenticalTo(g.Replaced(g.Target, s2).(Group).Target, s2) {
		t.Error("not replaced")
	}
	g = Group{Target: failingNoteConvertable{}}
	if !core.IsIdenticalTo(g.Replaced(s1, s2), g) {
		t.Error("should be same")
	}
}
