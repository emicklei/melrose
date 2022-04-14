package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestDecomposeSequence(t *testing.T) {
	s := core.MustParseSequence("1C 2D4+ =")
	ds := DecomposeSequence(s)
	t.Log(core.Storex(ds))
}
