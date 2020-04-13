package melrose

import "testing"

func TestJoin_Storex(t *testing.T) {
	l := MustParseSequence("A B")
	r := MustParseSequence("C D")

	if got, want := l.Join(r).Storex(), `join(sequence('A B'),sequence('C D'))`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestReverse_S(t *testing.T) {
	s := MustParseSequence("A B")

	if got, want := s.Reverse().S().String(), "B A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
