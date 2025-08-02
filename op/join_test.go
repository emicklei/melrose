package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestJoin_Storex(t *testing.T) {
	l := core.MustParseSequence("A B")
	r := core.MustParseSequence("C D")

	if got, want := (Join{Target: []core.Sequenceable{l, r}}).Storex(), `join(sequence('A B'),sequence('C D'))`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoin_S(t *testing.T) {
	l := core.MustParseSequence("A B")
	r := core.MustParseSequence("C D")
	j := Join{Target: []core.Sequenceable{l, r}}
	if got, want := j.S().Storex(), "sequence('A B C D')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	j = Join{Target: []core.Sequenceable{}}
	if got, want := j.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoin_Sequenceables(t *testing.T) {
	l := core.MustParseSequence("A B")
	j := Join{Target: []core.Sequenceable{l}}
	if got, want := len(j.Sequenceables()), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoin_Replaced(t *testing.T) {
	l := core.MustParseSequence("A B")
	r := core.MustParseSequence("C D")
	j := Join{Target: []core.Sequenceable{l, r}}
	if core.IsIdenticalTo(j, l) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(j.Replaced(l, r).(Join).Target[0], r) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(j.Replaced(j, r), r) {
		t.Error("should be replaced by r")
	}
}
