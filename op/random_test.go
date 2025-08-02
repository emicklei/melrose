package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestRandomInteger_Next(t *testing.T) {
	r := NewRandomInteger(core.On(1), core.On(100000))
	i1 := r.Next()
	if got, want := r.Value(), i1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := r.Value(), i1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r.Next()
	if got, want := r.Value(), i1; got == want {
		t.Errorf("got [%v:%T] do not want [%v:%T]", got, got, want, want)
	}
}

func TestRandomInteger_Storex(t *testing.T) {
	r := NewRandomInteger(core.On(1), core.On(10))
	if got, want := r.Storex(), "random(1,10)"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestRandomInteger_Next_Error(t *testing.T) {
	r := NewRandomInteger(core.On(10), core.On(1))
	if got, want := r.Next(), 10; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
