package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestOctave(t *testing.T) {
	s1 := melrose.MustParseSequence("C D E")
	s2 := melrose.MustParseSequence("F G A")
	o := Octave{Target: []melrose.Sequenceable{s1, s2}, Offset: melrose.On(1)}
	if got, want := o.S().Storex(), "sequence('C5 D5 E5 F5 G5 A5')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
