package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestOnBeat_Value(t *testing.T) {
	t.Skip()
	tl := new(melrose.TestLooper)
	tl.SetBIAB(4)
	tl.Tick() // off beat start
	i := melrose.NewInterval(melrose.On(0), melrose.On(12), melrose.On(2), melrose.RepeatFromTo)
	ob := NewOnBeat(melrose.On(4), i, tl)
	if got, want := ob.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tl.Tick() // 1
	if got, want := ob.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tl.Tick() // 2
	if got, want := ob.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tl.Tick() // 3
	if got, want := ob.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tl.Tick() // 4
	if got, want := ob.Value(), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tl.Tick()
	if got, want := ob.Value(), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
