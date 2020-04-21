package melrose

import (
	"testing"
)

func TestInterval_Value(t *testing.T) {
	i := NewInterval(On(0), On(1), On(1), RepeatFromTo)
	if got, want := i.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := i.Value(), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := i.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInterval_Value_Backwards(t *testing.T) {
	i := NewInterval(On(0), On(1), On(-1), RepeatFromTo)
	if got, want := i.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := i.Value(), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInterval_Value_Storex(t *testing.T) {
	i := NewInterval(On(0), On(1), On(-1), RepeatFromTo)
	if got, want := i.Storex(), "interval(0,1,-1,'repeat')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
