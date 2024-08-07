package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestScaleStepperNote(t *testing.T) {
	sc, _ := core.NewScale("E_")
	seq := core.MustParseSequence("A_ D")
	st := ScaleStepper{
		Scale:  core.On(sc),
		Count:  core.On(1),
		Target: seq,
	}
	if got, want := st.S().Storex(), "sequence('B_ E_5')"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
