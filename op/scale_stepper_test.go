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
	t.Log(st.S()) //  B_ E_5
}
