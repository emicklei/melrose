package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestMerge_EigthAndSixteenth(t *testing.T) {
	for _, each := range []struct {
		top, bottom, result string
	}{
		//{"F♯3 = F♯3", "16C♯5 16C♯ 16F♯ ⅛A", "sequence('(F♯3 16C♯5) 16C♯ 16F♯ (16= ⅛A) F#3"},
		{"8a 8a", "16d 16d 16d 16d", "sequence('(⅛A 16D) 16D (⅛A 16D) 16D')"},
		{"c", "d", "sequence('(C D)')"},
		{"c", "1d", "sequence('(C 1D)')"},
		{"=", "d", "sequence('(= D)')"},
		{"= = C (D E)", "= F = F F", "sequence('(= = =) (= = F) (C = =) (D E F) F')"},

		{"> e <", "f", "sequence('> (E F) <')"},
		{"> 8c 8= <", "16d 16d 16e 16e", "sequence('> (⅛C 16D) 16D (⅛= 16E) 16E <')"},
	} {
		s1 := core.MustParseSequence(each.top)
		s2 := core.MustParseSequence(each.bottom)
		m := Merge{Target: []core.Sequenceable{s1, s2}}
		if got, want := m.S().Storex(), each.result; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func TestReaderNoteAtDuration(t *testing.T) {
	r := sequenceReader{sequence: core.MustParseSequence("C D")}
	g, ok := r.noteStartingAt(0)
	if !ok {
		t.Fatal()
	}
	if got, want := g[0].Name, "C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	_, ok = r.noteStartingAt(0.15)
	if ok {
		t.Fatal()
	}
	g, ok = r.noteStartingAt(g[0].DurationFactor())
	if !ok {
		t.Fatal()
	}
	if got, want := g[0].Name, "D"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
