package op

import (
	"encoding/json"
	"fmt"
	"os"
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
		s := &core.Iterator{Target: []interface{}{c, d}}
		cn := core.MustParseNote("c")
		f := Fraction{Target: []core.Sequenceable{cn}, Parameter: core.On(1.0)}
		r := Replace{Target: f, From: cn, To: s}
		s.Next()
		if got, want := fmt.Sprintf("%v", r.S()), "1D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func TestReplace_JSON(t *testing.T) {
	c := core.MustParseSequence("C")
	d := core.MustParseSequence("D")
	p := Transpose{Target: c, Semitones: core.On(12)}
	r := Replace{Target: p, From: c, To: d}
	json.NewEncoder(os.Stdout).Encode(r)
}
