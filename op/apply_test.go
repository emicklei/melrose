package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

// f = apply(serial,octavemap,'1:-1,3:-1,1:0,2:0,3:0,1:1,2:1',parallel)
// f(a1)
func TestApply(t *testing.T) {
	serial := func(playable melrose.Sequenceable) melrose.Sequenceable {
		return playable
	}
	parallel := func(playable melrose.Sequenceable) melrose.Sequenceable {
		return playable
	}
	a := Apply{
		Target: []interface{}{serial, parallel},
	}
	c := melrose.MustParseChord("C")
	r := a.Call(c)
	t.Log(r.S().Notes)
}

func TestApplySerialInterface(t *testing.T) {
	serial := func(playables ...interface{}) interface{} {
		return playables[0]
	}
	a := Apply{
		Target: []interface{}{serial},
	}
	c := melrose.MustParseChord("C")
	r := a.Call(c)
	if r == nil {
		t.Fatal()
	}
	t.Log(r.S().Notes)
}

func TestApplyParallelInterface(t *testing.T) {
	parallel := func(value interface{}) interface{} {
		return value
	}
	a := Apply{
		Target: []interface{}{parallel},
	}
	c := melrose.MustParseChord("C")
	r := a.Call(c)
	if r == nil {
		t.Fatal()
	}
	t.Log(r.S().Notes)
}

func TestApplySerialParallelInterface(t *testing.T) {
	serial := func(playables ...interface{}) interface{} {
		return playables[0]
	}
	parallel := func(value interface{}) interface{} {
		return value
	}
	a := Apply{
		Target: []interface{}{serial, parallel},
	}
	c := melrose.MustParseChord("C")
	r := a.Call(c)
	if r == nil {
		t.Fatal()
	}
	t.Log(r.S().Notes)
}
