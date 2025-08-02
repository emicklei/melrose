package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestRepeat_Storex(t *testing.T) {
	s := core.MustParseSequence("C D")
	r := Repeat{Target: []core.Sequenceable{s}, Times: core.On(2)}
	if got, want := r.Storex(), "repeat(2,sequence('C D'))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestRepeat_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	r := Repeat{Target: []core.Sequenceable{s1}, Times: core.On(1)}
	if core.IsIdenticalTo(r, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(r.Replaced(s1, s2).(Repeat).Target[0], s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(r.Replaced(r, s2), s2) {
		t.Error("should be replaced by s2")
	}
}
