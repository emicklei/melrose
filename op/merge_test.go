package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestMerge_S(t *testing.T) {
	s1 := core.MustParseSequence("= = C (D E)")
	s2 := core.MustParseSequence("= F = F F")
	m := Merge{Target: []core.Sequenceable{s1, s2}}
	if got, want := m.S().Storex(), "sequence('= F C (D E F) F')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
