package core

import (
	"testing"
)

func TestScale_MajorC(t *testing.T) {
	s, _ := ParseScale("C")
	if got, want := s.S().Storex(), "sequence('C D E F G A B')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestScale_ChordAt(t *testing.T) {
	s, _ := ParseScale("C")
	for _, each := range []struct {
		step  int
		chord string
	}{
		{1, "chord('C')"},
		{2, "chord('D/m')"},
		{3, "chord('E/m')"},
		{4, "chord('F')"},
		{5, "chord('G')"},
		{6, "chord('A/m')"},
		{7, "chord('B')"},
	} {
		if got, want := s.ChordAt(each.step).Storex(), each.chord; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func TestScale_MinorC(t *testing.T) {
	s, _ := ParseScale("E/m")
	if got, want := s.S().Storex(), "sequence('E F G A B C5 D5')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestScale_MajorG(t *testing.T) {
	s, _ := ParseScale("G")
	if got, want := s.S().Storex(), "sequence('G A B C5 D5 E5 G_5')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
