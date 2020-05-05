package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestRepeat_Storex(t *testing.T) {
	s := melrose.MustParseSequence("C D")
	r := Repeat{Target: []melrose.Sequenceable{s}, Times: melrose.On(2)}
	if got, want := r.Storex(), "repeat(2,sequence('C D'))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
