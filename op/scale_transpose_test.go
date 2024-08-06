package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestScaleTransposeNote(t *testing.T) {

	s, _ := core.NewScale("E_")
	t.Log(core.N("A_").MIDI() / 12)
	t.Log(core.N("A_5").MIDI() % 12)

	i := s.IndexOfNote(core.N("A"))
	t.Log(i)
}
