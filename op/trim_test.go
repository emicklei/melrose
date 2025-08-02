package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestTrim_S(t *testing.T) {
	s := core.MustParseSequence("C D E F")
	tr := Trim{Target: s, Start: core.On(1), End: core.On(1)}
	if got, want := tr.S().Storex(), "sequence('D E')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	tr = Trim{Target: s, Start: core.On(-1), End: core.On(-1)}
	if got, want := tr.S().Storex(), "sequence('C D E F')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	tr = Trim{Target: s, Start: core.On(4), End: core.On(0)}
	if got, want := tr.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTrim_Storex(t *testing.T) {
	s := core.MustParseSequence("C")
	tr := Trim{Target: s, Start: core.On(1), End: core.On(1)}
	if got, want := tr.Storex(), "trim(1,1,sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	// tr = Trim{Target: failingNoteConvertable{}, Start: core.On(1), End: core.On(1)}
	// if got, want := tr.Storex(), ""; got != want {
	// 	t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	// }
}

func TestTrim_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	tr := Trim{Target: s1, Start: core.On(0), End: core.On(0)}
	if core.IsIdenticalTo(tr, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(tr.Replaced(s1, s2).(Trim).Target, s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(tr.Replaced(tr, s2), s2) {
		t.Error("should be replaced by s2")
	}
	tr = Trim{Target: tr, Start: core.On(0), End: core.On(0)}
	if !core.IsIdenticalTo(tr.Replaced(tr.Target, s2).(Trim).Target, s2) {
		t.Error("not replaced")
	}
	tr = Trim{Target: failingNoteConvertable{}, Start: core.On(0), End: core.On(0)}
	if !core.IsIdenticalTo(tr.Replaced(s1, s2), tr) {
		t.Error("should be same")
	}
}
