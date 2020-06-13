package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestAtIndex_S(t *testing.T) {
	s := melrose.MustParseSequence("C D")
	a := NewAtIndex(melrose.On(0), s)
	if got, want := len(a.S().Notes), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	a = NewAtIndex(melrose.On(1), s)
	if got, want := len(a.S().Notes), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	a = NewAtIndex(melrose.On(3), s)
	if got, want := len(a.S().Notes), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := a.Storex(), "at(3,sequence('C D'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
