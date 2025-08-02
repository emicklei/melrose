package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestSerial_S(t *testing.T) {
	s := core.MustParseSequence("(C D) E")
	g := Serial{Target: []core.Sequenceable{s}}
	if got, want := g.S().Storex(), "sequence('C D E')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestSerial_Storex(t *testing.T) {
	s := core.MustParseSequence("C D E")
	g := Serial{Target: []core.Sequenceable{s}}
	if got, want := g.Storex(), "ungroup(sequence('C D E'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	// g = Serial{Target: []core.Sequenceable{failingNoteConvertable{}}}
	// if got, want := g.Storex(), ""; got != want {
	// 	t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	// }
}

func TestSerial_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	g := Serial{Target: []core.Sequenceable{s1}}
	if core.IsIdenticalTo(g, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(g.Replaced(s1, s2).(Serial).Target[0], s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(g.Replaced(g, s2), s2) {
		t.Error("should be replaced by s2")
	}
}
