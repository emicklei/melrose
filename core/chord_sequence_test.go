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

func TestChordSequence_Replaced(t *testing.T) {
	s1 := MustParseChordSequence("C D")
	s2 := MustParseChordSequence("E F")
	s3 := MustParseChordSequence("C D") // Identical to s1

	// Case 1: p is identical to from
	replaced1 := s1.Replaced(s3, s2)
	if !IsIdenticalTo(replaced1, s2) {
		t.Errorf("Expected s1.Replaced(s3, s2) to be s2, got %v", replaced1)
	}

	// Case 2: p is not identical to from
	replaced2 := s1.Replaced(s2, s3)
	if !IsIdenticalTo(replaced2, s1) {
		t.Errorf("Expected s1.Replaced(s2, s3) to be s1, got %v", replaced2)
	}
}
