package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestRotate_S(t *testing.T) {
	s := core.MustParseSequence("C D E")
	r := Rotate{Target: s, Times: core.On(1)}
	if got, want := r.S().Storex(), "sequence('E C D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestRotate_Storex(t *testing.T) {
	s := core.MustParseSequence("C D E")
	r := Rotate{Target: s, Times: core.On(1)}
	if got, want := r.Storex(), "rotate(1,sequence('C D E'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
