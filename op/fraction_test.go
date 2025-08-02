package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestFraction_S(t *testing.T) {
	f := NewFraction(core.On(0.75), core.InList(core.N("c")))
	if got, want := f.S().Storex(), "sequence('2.C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	f = NewFraction(core.On(0.375), core.InList(core.N("c")))
	if got, want := f.S().Storex(), "sequence('.C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	f = NewFraction(core.On(2), core.InList(core.N("c")))
	if got, want := f.S().Storex(), "sequence('2C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestFraction_Storex(t *testing.T) {
	f := NewFraction(core.On(2), core.InList(core.N("c")))
	if got, want := f.Storex(), "fraction(2,note('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestFraction_ToNote(t *testing.T) {
	f := NewFraction(core.On(2), core.InList(core.N("c")))
	n, err := f.ToNote()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := n.Storex(), "note('2C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	f = NewFraction(core.On(2), []core.Sequenceable{})
	_, err = f.ToNote()
	if err == nil {
		t.Fatal("error expected")
	}
	f = NewFraction(core.On(2), core.InList(core.S("c d")))
	_, err = f.ToNote()
	if err == nil {
		t.Fatal("error expected")
	}
	f = NewFraction(core.On(2), core.InList(failingNoteConvertable{}))
	_, err = f.ToNote()
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestFraction_Replaced(t *testing.T) {
	f := NewFraction(core.On(2), core.InList(core.N("c")))
	if core.IsIdenticalTo(f, f.Target[0]) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(f.Replaced(f.Target[0], core.EmptySequence).(Fraction).Target[0], core.EmptySequence) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(f.Replaced(f, core.EmptySequence), core.EmptySequence) {
		t.Error("should be replaced by empty")
	}
}
