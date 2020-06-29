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
func TestScale_MinorC(t *testing.T) {
	s, _ := ParseScale("E/m")
	if got, want := s.S().Storex(), "sequence('E F G A B C5 D5')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestScale_MajorG(t *testing.T) {
	s, _ := ParseScale("G")
	if got, want := s.S().Storex(), "sequence('G A B C5 D5 E5 G♭5')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestScale_TwoOctaves(t *testing.T) {
	s, _ := NewScale(2, "e")
	if got, want := s.S().Storex(), "sequence('E G♭ A♭ A B D♭5 E♭5 E5 G♭5 A♭5 A5 B5 D♭6 E♭6')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
