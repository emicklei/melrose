package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestOctave(t *testing.T) {
	s1 := core.MustParseSequence("C D E~8E")
	s2 := core.MustParseSequence("F G A")
	o := Octave{Target: []core.Sequenceable{s1, s2}, Offset: core.On(1)}
	if got, want := o.S().Storex(), "sequence('C5 D5 E5~8E5 F5 G5 A5')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
