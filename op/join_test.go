package op

import (
	"github.com/emicklei/melrose/core"
	"testing"
)

func TestJoin_Storex(t *testing.T) {
	l := core.MustParseSequence("A B")
	r := core.MustParseSequence("C D")

	if got, want := (Join{Target: []core.Sequenceable{l, r}}).Storex(), `join(sequence('A B'),sequence('C D'))`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
