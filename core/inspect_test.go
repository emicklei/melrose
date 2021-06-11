package core

import "testing"

func TestInspection_Markdown_Chord(t *testing.T) {
	c := MustParseChord("b3#/m/1")
	i := NewInspect(testContext(), "", c)
	t.Log(i.Markdown())
}

func TestInspection_Markdown_ChordSequence(t *testing.T) {
	c := MustParseChordSequence("b3#/m/1 = c3/2")
	i := NewInspect(testContext(), "", c)
	t.Log(i.Markdown())
}
