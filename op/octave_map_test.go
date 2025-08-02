package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestOctaveMapper_S(t *testing.T) {
	o := NewOctaveMap(core.MustParseSequence("C (D E) F"), "1:-1,2:1")
	if got, want := o.S().Storex(), "sequence('C3 (D5 E5)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}

}

func TestOctaveMapper_parseIndices(t *testing.T) {
	m := parseIndexOffsets("1:-1,3:-1,1:0,2:0,3:0,1:1,2:1")
	if got, want := m[0].from, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(m), 7; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestOctaveMap_Storex(t *testing.T) {
	o := NewOctaveMap(core.MustParseSequence("C (D E) F"), "1:-1,2:1")
	if got, want := o.Storex(), "octavemap('1:-1,2:1',sequence('C (D E) F'))"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	// o = NewOctaveMap(failingNoteConvertable{}, "1:-1,2:1")
	// if got, want := o.Storex(), ""; got != want {
	// 	t.Errorf("got [%v] want [%v]", got, want)
	// }
}

func TestOctaveMap_Notes(t *testing.T) {
	o := NewOctaveMap(core.MustParseSequence("C (D E) F"), "1:-1,2:1,4:1")
	if got, want := len(o.Notes()), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	o = NewOctaveMap(core.MustParseSequence("C (D E) F"), "1:0")
	if got, want := o.Notes()[0][0].Octave, 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestOctaveMap_Replaced(t *testing.T) {
	s1 := core.MustParseSequence("C")
	s2 := core.MustParseSequence("D")
	o := NewOctaveMap(s1, "1:1")
	if core.IsIdenticalTo(o, s1) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(o.Replaced(s1, s2).(OctaveMap).Target, s2) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(o.Replaced(o, s2), s2) {
		t.Error("should be replaced by s2")
	}
	o = NewOctaveMap(o, "1:1")
	if !core.IsIdenticalTo(o.Replaced(o.Target, s2).(OctaveMap).Target, s2) {
		t.Error("not replaced")
	}
	o = NewOctaveMap(failingNoteConvertable{}, "1:1")
	if !core.IsIdenticalTo(o.Replaced(s1, s2), o) {
		t.Error("should be same")
	}
}
