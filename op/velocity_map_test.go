package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestVelocityMap_S(t *testing.T) {
	o := NewVelocityMap(core.MustParseSequence("C (D E) F"), "1:30,2:60")
	if got, want := o.S().Storex(), "sequence('C--- (D+ E+)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}

}
