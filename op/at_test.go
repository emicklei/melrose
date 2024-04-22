package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestAtIndex_S(t *testing.T) {
	s := core.MustParseSequence("C D")
	a := NewAtIndex(core.On(0), s)
	if got, want := len(a.S().Notes), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	a = NewAtIndex(core.On(1), s)
	if got, want := len(a.S().Notes), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	a = NewAtIndex(core.On(3), s)
	if got, want := len(a.S().Notes), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := a.Storex(), "at(3,sequence('C D'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	b := NewAtIndex(core.On(1), s).S()
	if got, want := b.Storex(), "sequence('C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	b2 := NewAtIndex(core.On(2), s).S()
	if got, want := b2.Storex(), "sequence('D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
