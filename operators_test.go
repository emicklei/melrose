package melrose

import (
	"testing"
)

func TestReverse_S(t *testing.T) {
	s := MustParseSequence("A B")

	if got, want := s.Reverse().S().String(), "B A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestJoin_S(t *testing.T) {
	l := MustParseSequence("A B")
	r := MustParseSequence("C D")

	if got, want := l.Join(r).S().String(), "A B C D"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
