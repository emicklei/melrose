package core

import (
	"testing"
)

func TestNewChordProgression(t *testing.T) {
	p := NewChordProgression(On("C"), On("II V I"))
	if got, want := p.S().Storex(), "sequence('(D F A) (G B D5) (C E G)')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
