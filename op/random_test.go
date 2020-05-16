package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestRandomInteger_Next(t *testing.T) {
	r := NewRandomInteger(melrose.On(1), melrose.On(100000))
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
