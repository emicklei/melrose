package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestOctaveMapper_S(t *testing.T) {
	o := NewOctaveMap(core.MustParseSequence("C (D E) F"), "1:-1,2:1")
	if got, want := o.S().Storex(), "sequence('C3 (D5 E5)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}

}

func TestOctaveMapper_parseIndices(t *testing.T) {
	m := parseIndexOffsets("1:-1,3:-1,1:0,2:0,3:0,1:1,2:1")
	if got, want := m[0].from, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(m), 7; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
