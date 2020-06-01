package op

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/emicklei/melrose"
)

func TestReplace_Operators(t *testing.T) {
	c := melrose.MustParseSequence("C")
	d := melrose.MustParseSequence("D")

	{
		r := Replace{Target: c, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		p := Pitch{Target: c, Semitones: melrose.On(12)}
		r := Replace{Target: p, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D5"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		rp := Repeat{Target: []melrose.Sequenceable{d, c}, Times: melrose.On(1)}
		r := Replace{Target: rp, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		j := Join{Target: []melrose.Sequenceable{c, d}}
		r := Replace{Target: j, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		s := SequenceMapper{Target: c, Indices: [][]int{[]int{1}}}
		r := Replace{Target: s, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
	{
		s := Octave{Target: []melrose.Sequenceable{c}, Offset: melrose.On(1)}
		r := Replace{Target: s, From: c, To: d}
		if got, want := fmt.Sprintf("%v", r.S()), "D5"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	}
}

func TestReplace_JSON(t *testing.T) {
	c := melrose.MustParseSequence("C")
	d := melrose.MustParseSequence("D")
	p := Pitch{Target: c, Semitones: melrose.On(12)}
	r := Replace{Target: p, From: c, To: d}
	json.NewEncoder(os.Stdout).Encode(r)
}
