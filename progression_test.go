package melrose

import (
	"testing"
)

func TestParseProgression(t *testing.T) {
	empty, err := ParseProgression("")
	check(t, err)
	if got, want := len(empty.Chords), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	two, err := ParseProgression("C D")
	check(t, err)
	if got, want := len(two.Chords), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := two.S().Storex(), "sequence('(C E G) (D G♭ A)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestParseProgression_ParallelChords(t *testing.T) {
	par := MustParseProgression("(E F)")
	if got, want := par.S().Storex(), "sequence('(E A♭ B F A C5)')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestParseProgression_Storex(t *testing.T) {
	par := MustParseProgression("A (E F) =")
	if got, want := par.Storex(), "progression('A (E F) =')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
