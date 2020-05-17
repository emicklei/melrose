package dsl

import (
	"testing"
)

func Test_variable_Sub(t *testing.T) {
	s := NewVariableStore()
	s.Put("a", 1)
	a := s.getVariable("a")
	if got, want := a.dispatchAdd(1), 2; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := a.dispatchAdd(a), 2; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := a.dispatchSub(1), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := a.dispatchSub(a), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := a.dispatchSubFrom(2), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := a.dispatchSubFrom(a), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
