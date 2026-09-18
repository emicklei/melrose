package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestTranspose_S(t *testing.T) {
	s := core.MustParseSequence("C")
	tr := Transpose{Target: []core.Sequenceable{s}, Semitones: core.On(1)}
	if got, want := tr.S().Storex(), "sequence('D_')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTranspose_Storex(t *testing.T) {
	s := core.MustParseSequence("C")
	tr := Transpose{Target: []core.Sequenceable{s}, Semitones: core.On(1)}
	if got, want := tr.Storex(), "transpose(1,sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTranspose_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	tr := Transpose{Target: []core.Sequenceable{s1}, Semitones: core.On(1)}
	if core.IsIdenticalTo(tr, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(tr.Replaced(s1, s2).(Transpose).Target[0], s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(tr.Replaced(tr, s2), s2) {
		t.Error("should be replaced by s2")
	}
	tr = Transpose{Target: []core.Sequenceable{tr}, Semitones: core.On(1)}
	if !core.IsIdenticalTo(tr.Replaced(tr.Target[0], s2).(Transpose).Target[0], s2) {
		t.Error("not replaced")
	}
	tr = Transpose{Target: []core.Sequenceable{failingNoteConvertable{}}, Semitones: core.On(1)}
	if !core.IsIdenticalTo(tr.Replaced(s1, s2), tr) {
		t.Error("should be same")
	}
}
