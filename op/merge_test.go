package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestMerge(t *testing.T) {
	for _, each := range []struct {
		top, bottom, result string
	}{
		{"8F#3 8= 8F#3", "16C#5 16C# 16F# 8A 16A", "sequence('(8F#3 16C#5) 16C# 16F# (16= 8A) (8F#3 16=) 16A')"},
		{"8a 8a", "16d 16d 16d 16d", "sequence('(8A 16D) 16D (8A 16D) 16D')"},
		{"c", "d", "sequence('(C D)')"},
		{"c", "1d", "sequence('(C 1D) 2.=')"},
		{"=", "d", "sequence('D')"},
		{"= = C (D E)", "= F = F F", "sequence('= F C (D E F) F')"},
		{"> e <", "f", "sequence('> (E F) <')"},
		{"> 8c 8= <", "16d 16d 16e 16e", "sequence('> (8C 16D) 16D 16E 16E <')"},
		{"> C <", "> D <", "sequence('> (C D) <')"},
		{"> C <", "> D <", "sequence('> (C D) <')"},
	} {
		s1 := core.MustParseSequence(each.top)
		s2 := core.MustParseSequence(each.bottom)
		m := Merge{Target: []core.Sequenceable{s1, s2}}
		if got, want := m.S().Storex(), each.result; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func Test_compactGroup(t *testing.T) {
	for _, each := range []struct {
		in, out string
	}{
		{"(C)", "sequence('C')"},
		{"(C =)", "sequence('C')"},
		{"(16= 16C#)", "sequence('16C#')"},
		{"(16= 8A)", "sequence('(16= 8A)')"},
	} {
		g := core.MustParseSequence(each.in)
		c := core.Sequence{Notes: [][]core.Note{compactGroup(g.Notes[0])}}
		if got, want := c.Storex(), each.out; got != want {
			t.Errorf("got [%v:%T] want [%v:%T] from [%s]", got, got, want, want, g.Storex())
		}
	}
}

func TestMerge_Storex(t *testing.T) {
	s1 := core.MustParseSequence("C")
	m := Merge{Target: []core.Sequenceable{s1}}
	if got, want := m.Storex(), "merge(sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestMerge_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	m := Merge{Target: []core.Sequenceable{s1}}
	if core.IsIdenticalTo(m, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(m.Replaced(s1, s2).(Join).Target[0], s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(m.Replaced(m, s2), s2) {
		t.Error("should be replaced by s2")
	}
}
