package core

import (
	"testing"
)

func TestPrint_Interval(t *testing.T) {
	i := NewInterval(On(1), On(2), On(1), RepeatFromTo)
	n := Nexter{Target: i}
	var w Sequenceable = Print{Target: n}
	if _, ok := w.(Sequenceable); !ok {
		t.Fail()
	}
}
