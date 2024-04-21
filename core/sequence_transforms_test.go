package core

import "testing"

func TestTransforms(t *testing.T) {
	s, e := ParseSequence("16F2-- = 16G2-- = 16A2-- = 16G2-- = 16F2-- = 16G2-- = 8A2-- = 16B2-- = 16C3-- = 16B2--- = 8A2-- 8.= 8G2-- 8.= 8A2-- 8.= 8B2-- 8.= 8C3-- 2= 8C3-- 16= 8D3++ 16= 8E3+++ 16= 16D3++ 16= 16C3+++ 1.= 16E2-- 8= 8F2++ 16= 16G2 8= 16A2++ 8= 16B2+++ 16= 8A2++ 16= 16G2++ 8= 16F2++")
	if e != nil {
		t.Fatal(e)
	}
	t.Log(s.NoFractions().Storex())
	t.Log(s.NoFractions().NoDynamics().Storex())
	t.Log(s.NoFractions().NoDynamics().NoRests().Storex())
}
