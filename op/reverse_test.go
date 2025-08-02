package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestReverse_S(t *testing.T) {
	s := core.MustParseSequence("A B")

	if got, want := (Reverse{Target: s}).S().String(), "B A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestReverse_Storex(t *testing.T) {
	s := core.MustParseSequence("A B")
	r := Reverse{Target: s}
	if got, want := r.Storex(), "reverse(sequence('A B'))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	// r = Reverse{Target: failingNoteConvertable{}}
	// if got, want := r.Storex(), ""; got != want {
	// 	t.Errorf("got [%v] want [%v]", got, want)
	// }
}

func TestReverse_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	r := Reverse{Target: s1}
	if core.IsIdenticalTo(r, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(r.Replaced(s1, s2).(Reverse).Target, s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(r.Replaced(r, s2), s2) {
		t.Error("should be replaced by s2")
	}
	r = Reverse{Target: r}
	if !core.IsIdenticalTo(r.Replaced(r.Target, s2).(Reverse).Target, s2) {
		t.Error("not replaced")
	}
	r = Reverse{Target: failingNoteConvertable{}}
	if !core.IsIdenticalTo(r.Replaced(s1, s2), r) {
		t.Error("should be same")
	}
}
