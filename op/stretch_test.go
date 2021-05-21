package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestStretch_S(t *testing.T) {
	s := core.MustParseSequence("(c d) 2e 8.f")
	st := s.Stretched(2.0)
	if got, want := core.Storex(st), "sequence('(2C 2D) 1E .F')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
