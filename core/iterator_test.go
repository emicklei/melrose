package core

import (
	"testing"
)

func TestIterator_Value(t *testing.T) {
	var l []interface{} = []interface{}{"C", "D"}
	i := &Iterator{Target: l}
	if got, want := i.Value(), "C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	i.Next()
	if got, want := i.Value(), "D"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	i.Next()
	if got, want := i.Value(), "C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
