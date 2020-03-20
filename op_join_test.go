package melrose

import "testing"

func TestJoin_S(t *testing.T) {
	l := MustParseSequence("A B")
	r := MustParseSequence("C D")

	if got, want := l.Join(r).S().String(), "A B C D"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoin_Storex(t *testing.T) {
	l := MustParseSequence("A B")
	r := MustParseSequence("C D")

	if got, want := l.Join(r).Storex(), `join(seq("A B"),seq("C D"))`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
