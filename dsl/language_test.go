package dsl

import (
	"testing"
)

func TestNote(t *testing.T) {
	r := eval(t, "note('c')")
	checkStorex(t, r, "note('C')")
}

func TestChord(t *testing.T) {
	r := eval(t, "chord('C#/m')")
	checkStorex(t, r, "chord('C♯/m')")
}
