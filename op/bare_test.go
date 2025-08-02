package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestBare_S(t *testing.T) {
	s := core.MustParseSequence("C D E")
	b := Bare{Target: []core.Sequenceable{s}}
	if got, want := b.S().Storex(), "sequence('C D E')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestBare_Storex(t *testing.T) {
	s := core.MustParseSequence("C D E")
	b := Bare{Target: []core.Sequenceable{s}}
	if got, want := b.Storex(), "bare(sequence('C D E'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestBare_Replaced(t *testing.T) {
	s := core.MustParseSequence("C")
	d := core.MustParseSequence("D")
	b := Bare{Target: []core.Sequenceable{s}}
	if core.IsIdenticalTo(b, s) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(b.Replaced(s, d).(Bare).Target[0], d) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(b.Replaced(d, s), b) {
		t.Error("should be same")
	}
	if !core.IsIdenticalTo(b.Replaced(b, s), s) {
		t.Error("should be replaced by s")
	}
}
