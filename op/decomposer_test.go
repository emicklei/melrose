package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestDecomposeSequence(t *testing.T) {
	s := core.MustParseSequence("1C# 2D_4+ = (2e-- 2e5--)")
	ds := DecomposeSequence(s)
	st := core.Storex(ds)
	if st != "dynamicmap('2:+,4:--',fractionmap('1 2 4 2',sequence('C# D_ = (E E5)')))" {
		t.Errorf("expected 'sequence(...)' got %s", st)
	}
}
