package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestTranspose_S(t *testing.T) {
	s := core.MustParseSequence("C")
	tr := Transpose{Target: s, Semitones: core.On(1)}
	if got, want := tr.S().Storex(), "sequence('D_')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTranspose_Storex(t *testing.T) {
	s := core.MustParseSequence("C")
	tr := Transpose{Target: s, Semitones: core.On(1)}
	if got, want := tr.Storex(), "transpose(1,sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTranspose_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	tr := Transpose{Target: s1, Semitones: core.On(1)}
	if core.IsIdenticalTo(tr, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(tr.Replaced(s1, s2).(Transpose).Target, s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(tr.Replaced(tr, s2), s2) {
		t.Error("should be replaced by s2")
	}
	tr = Transpose{Target: tr, Semitones: core.On(1)}
	if !core.IsIdenticalTo(tr.Replaced(tr.Target, s2).(Transpose).Target, s2) {
		t.Error("not replaced")
	}
	tr = Transpose{Target: failingNoteConvertable{}, Semitones: core.On(1)}
	if !core.IsIdenticalTo(tr.Replaced(s1, s2), tr) {
		t.Error("should be same")
	}
}
