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
