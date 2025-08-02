package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestTransposeMap_S(t *testing.T) {
	s := core.MustParseSequence("C D E")
	tr := NewTransposeMap(s, "1:1,3:2")
	if got, want := tr.S().Storex(), "sequence('D_ G_')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTransposeMap_Notes(t *testing.T) {
	s := core.MustParseSequence("C D E")
	tr := NewTransposeMap(s, "1:1,4:2")
	if got, want := len(tr.Notes()), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tr = NewTransposeMap(s, "1:0")
	if got, want := tr.Notes()[0][0].MIDI(), 60; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTransposeMap_Storex(t *testing.T) {
	s := core.MustParseSequence("C")
	tr := NewTransposeMap(s, "1:1")
	if got, want := tr.Storex(), "transposemap('1:1',sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTransposeMap_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	tr := NewTransposeMap(s1, "1:1")
	if core.IsIdenticalTo(tr, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(tr.Replaced(s1, s2).(TransposeMap).Target, s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(tr.Replaced(tr, s2), s2) {
		t.Error("should be replaced by s2")
	}
	tr = NewTransposeMap(tr, "1:1")
	if !core.IsIdenticalTo(tr.Replaced(tr.Target, s2).(TransposeMap).Target, s2) {
		t.Error("not replaced")
	}
	tr = NewTransposeMap(failingNoteConvertable{}, "1:1")
	if !core.IsIdenticalTo(tr.Replaced(s1, s2), tr) {
		t.Error("should be same")
	}
}
