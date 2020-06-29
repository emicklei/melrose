package op

import (
	"github.com/emicklei/melrose/core"
	"testing"
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
