package op

import (
	"fmt"
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestReplace_Operators(t *testing.T) {
	c := core.MustParseSequence("C")
	d := core.MustParseSequence("D")

	{
		r := Replace{Target: c, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		p := Transpose{Target: c, Semitones: core.On(12)}
		r := Replace{Target: p, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D5"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		rp := Repeat{Target: []core.Sequenceable{d, c}, Times: core.On(1)}
		r := Replace{Target: rp, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		j := Join{Target: []core.Sequenceable{c, d}}
		r := Replace{Target: j, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		s := Resequencer{Target: c, Indices: [][]int{[]int{1}}}
		r := Replace{Target: s, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		s := Octave{Target: []core.Sequenceable{c}, Offset: core.On(1)}
		r := Replace{Target: s, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D5"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		s := &core.Iterator{Target: []any{c, d}}
		cn := core.MustParseNote("c")
		f := Fraction{Target: []core.Sequenceable{cn}, Parameter: core.On(1.0)}
		r := Replace{Target: f, From: cn, To: s}
		s.Next()
		if got, want := fmt.Sprintf("%v", r.S()), "1D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func TestReplace_Storex(t *testing.T) {
	c := core.MustParseSequence("C")
	d := core.MustParseSequence("D")
	r := Replace{Target: c, From: c, To: d}
	if got, want := r.Storex(), "replace(sequence('C'),sequence('C'),sequence('D'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r = Replace{Target: c, From: c, To: core.MustParseSequence("1")}
	if got, want := r.Storex(), "replace(sequence('C'),sequence('C'),sequence('1'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r = Replace{Target: c, From: core.MustParseSequence("1"), To: d}
	if got, want := r.Storex(), "replace(sequence('C'),sequence('1'),sequence('D'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	r = Replace{Target: core.MustParseSequence("1"), From: c, To: d}
	if got, want := r.Storex(), "replace(sequence('1'),sequence('C'),sequence('D'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestReplace_Replaced(t *testing.T) {
	c := core.MustParseSequence("C")
	d := core.MustParseSequence("D")
	r := Replace{Target: c, From: c, To: d}
	if core.IsIdenticalTo(r, c) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(r.Replaced(r, d), d) {
		t.Error("should be replaced by d")
	}
	r = Replace{Target: r, From: c, To: d}
	if !core.IsIdenticalTo(r.Replaced(r.Target, d).(Replace).Target, d) {
		t.Error("not replaced")
	}
	r = Replace{Target: failingNoteConvertable{}, From: c, To: d}
	if !core.IsIdenticalTo(r.Replaced(c, d), r) {
		t.Error("should be same")
	}
}
