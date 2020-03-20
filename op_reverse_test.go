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
