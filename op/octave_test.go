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

func TestOctave_Storex(t *testing.T) {
	s1 := core.MustParseSequence("C D E~8E")
	o := Octave{Target: []core.Sequenceable{s1}, Offset: core.On(1)}
	if got, want := o.Storex(), "octave(1,sequence('C D E~8E'))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestOctave_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	o := Octave{Target: []core.Sequenceable{s1}, Offset: core.On(1)}
	if core.IsIdenticalTo(o, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(o.Replaced(s1, s2).(Octave).Target[0], s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(o.Replaced(o, s2), s2) {
		t.Error("should be replaced by s2")
	}
}
