package core

import (
	"testing"
)

func TestParseChordSequence(t *testing.T) {
	empty, err := ParseChordSequence("")
	check(t, err)
	if got, want := len(empty.Chords), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	two, err := ParseChordSequence("C D")
	check(t, err)
	if got, want := len(two.Chords), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := two.S().Storex(), "sequence('(C E G) (D G_ A)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestParseChordSequence_ParallelChords(t *testing.T) {
	par := MustParseChordSequence("(E F)")
	if got, want := par.S().Storex(), "sequence('(E A_ B F A C5)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestParseChordSequence_Storex(t *testing.T) {
	par := MustParseChordSequence("A (E F) =")
	if got, want := par.Storex(), "chordsequence('A (E F) =')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
