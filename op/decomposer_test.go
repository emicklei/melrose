package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestDecomposeSequence(t *testing.T) {
	s := core.MustParseSequence("1C# 2D_4+ = (2e-- 2e5--)")
	ds := DecomposeSequence(s)
	t.Log(core.Storex(ds))
}
