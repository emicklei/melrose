package core

import (
	"testing"
)

func TestTrack_Add(t *testing.T) {
	tr := NewTrack("test", 1)
	s1 := MustParseSequence("C D E F")
	s2 := MustParseSequence("G A B C5")
	tr.Add(s1)
	tr.Add(s2)
	if got, want := tr.Content[1], s1; storex(got) != storex(want) {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := tr.Content[2], s2; storex(got) != storex(want) {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
