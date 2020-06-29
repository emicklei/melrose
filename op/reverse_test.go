package op

import (
	"github.com/emicklei/melrose/core"
	"testing"
)

func TestReverse_S(t *testing.T) {
	s := core.MustParseSequence("A B")

	if got, want := (Reverse{Target: s}).S().String(), "B A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
