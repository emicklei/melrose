package core

import (
	"testing"
)

func TestNearest(t *testing.T) {
	if got, want := nearest(5, 3), int64(6); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := nearest(539, 125), int64(500); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := nearest(8, 10), int64(10); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestNotePeriod_Quantized(t *testing.T) {
	bpm := 100.0
	p := NotePeriod{startMs: 0, endMs: 100, number: 65, velocity: Normal}
	w := WholeNoteDuration(bpm) / 4
	p = p.Quantized(bpm)
	n := p.Note(bpm)
	t.Logf("%v %#v %v", w, p, n)
	{
		p := NotePeriod{startMs: 0, endMs: 750, number: 65, velocity: Normal}
		w := WholeNoteDuration(bpm) / 4
		p = p.Quantized(bpm)
		n := p.Note(bpm)
		t.Logf("%v %#v %v", w, p, n)
	}
	{
		p := NotePeriod{startMs: 0, endMs: 800, number: 65, velocity: Normal}
		w := WholeNoteDuration(bpm) / 4
		p = p.Quantized(bpm)
		n := p.Note(bpm)
		t.Logf("%v %#v %v", w, p, n)
	}
}
