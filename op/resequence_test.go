package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestResequence_S(t *testing.T) {
	s := core.MustParseSequence("C D E")
	r := NewResequencer(s, core.On("1 3 2"))
	if got, want := r.S().Storex(), "sequence('C E D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r = NewResequencer(s, core.On(""))
	if got, want := r.S().Storex(), "sequence('C D E')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r = NewResequencer(s, core.On("1 4 2"))
	if got, want := r.S().Storex(), "sequence('C D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r = NewResequencer(s, nil)
	if got, want := r.S().Storex(), "sequence('C D E')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestResequence_Storex(t *testing.T) {
	s := core.MustParseSequence("C D E")
	r := NewResequencer(s, core.On("1 3 2"))
	if got, want := r.Storex(), "resequence('1 3 2',sequence('C D E'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	// r = NewResequencer(failingNoteConvertable{}, core.On("1"))
	// if got, want := r.Storex(), "?"; got != want {
	// 	t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	// }
	// r = NewResequencer(s, core.On(1))
	// if got, want := r.Storex(), "?"; got != want {
	// 	t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	// }
}

func TestResequence_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	r := NewResequencer(s1, core.On("1"))
	if core.IsIdenticalTo(r, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(r.Replaced(s1, s2).(Resequencer).Target, s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(r.Replaced(r, s2), s2) {
		t.Error("should be replaced by s2")
	}
	r = NewResequencer(r, core.On("1"))
	if !core.IsIdenticalTo(r.Replaced(r.Target, s2).(Resequencer).Target, s2) {
		t.Error("not replaced")
	}
	r = NewResequencer(failingNoteConvertable{}, core.On("1"))
	if !core.IsIdenticalTo(r.Replaced(s1, s2), r) {
		t.Error("should be same")
	}
}
