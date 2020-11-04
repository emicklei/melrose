package core

import (
	"testing"
)

func TestTrack_Add(t *testing.T) {
	tr := NewTrack("test", 1)
	s1 := MustParseSequence("C D E F")
	tr.Add(NewSequenceOnTrack(On(1), s1))
	if got, want := Storex(tr), "track('test',1,onbar(1,sequence('C D E F')))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
