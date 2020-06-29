package op

import (
	"github.com/emicklei/melrose/core"
	"testing"
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
}
