package op

import (
	"testing"

	. "github.com/emicklei/melrose"
)

func TestReverse_S(t *testing.T) {
	s := MustParseSequence("A B")

	if got, want := (Reverse{Target: s}).S().String(), "B A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
