package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestJoin_Storex(t *testing.T) {
	l := melrose.MustParseSequence("A B")
	r := melrose.MustParseSequence("C D")

	if got, want := (Join{Target: []melrose.Sequenceable{l, r}}).Storex(), `join(sequence('A B'),sequence('C D'))`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
