package core

import (
	"testing"
)

func TestPrint_Interval(t *testing.T) {
	i := NewInterval(On(1), On(2), On(1), RepeatFromTo)
	n := Nexter{Target: i}
	var w Sequenceable = Print{Target: n}
	if got, want := Storex(w), "print(next(interval(1,2,1,'repeat')))"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
