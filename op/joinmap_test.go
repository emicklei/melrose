package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestJoinMapper_S(t *testing.T) {
	j := Join{Target: []melrose.Sequenceable{melrose.MustParseNote("c"), melrose.MustParseNote("d")}}
	m := NewJoinMapper(melrose.On(j), "(1 2) 1")
	if got, want := m.S().Storex(), "sequence('(C D) C')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
