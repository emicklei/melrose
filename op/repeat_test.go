package op

import (
	"github.com/emicklei/melrose/core"
	"testing"
)

func TestRepeat_Storex(t *testing.T) {
	s := core.MustParseSequence("C D")
	r := Repeat{Target: []core.Sequenceable{s}, Times: core.On(2)}
	if got, want := r.Storex(), "repeat(2,sequence('C D'))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
