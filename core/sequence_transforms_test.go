package core

import "testing"

func TestTransforms(t *testing.T) {
	s, e := ParseSequence("16E2+++ 8= 16F2++++ 16= 16G2+++")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := len(s.NoFractions().NoDynamics().NoRests().Notes), 3; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
