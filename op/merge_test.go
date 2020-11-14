package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestMerge(t *testing.T) {
	for _, each := range []struct {
		top, bottom, result string
	}{
		{"⅛F♯3 ⅛= ⅛F♯3", "16C♯5 16C♯ 16F♯ ⅛A 16A", "sequence('(⅛F♯3 16C♯5) 16C♯ 16F♯ (16= ⅛A) (⅛F♯3 16=) 16A')"},
		{"8a 8a", "16d 16d 16d 16d", "sequence('(⅛A 16D) 16D (⅛A 16D) 16D')"},
		{"c", "d", "sequence('(C D)')"},
		{"c", "1d", "sequence('(C 1D) ½.=')"},
		{"=", "d", "sequence('D')"},
		{"= = C (D E)", "= F = F F", "sequence('= F C (D E F) F')"},
		{"> e <", "f", "sequence('> (E F) <')"},
		{"> 8c 8= <", "16d 16d 16e 16e", "sequence('> (⅛C 16D) 16D 16E 16E <')"},
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
		{"(16= 16C♯)", "sequence('16C♯')"},
		{"(16= ⅛A)", "sequence('(16= ⅛A)')"},
	} {
		g := core.MustParseSequence(each.in)
		c := core.Sequence{Notes: [][]core.Note{compactGroup(g.Notes[0])}}
		if got, want := c.Storex(), each.out; got != want {
			t.Errorf("got [%v:%T] want [%v:%T] from [%s]", got, got, want, want, g.Storex())
		}
	}
}
