package op

import (
	"testing"

	"github.com/emicklei/melrose"
)

// f = pipeline(serial,octavemap,'1:-1,3:-1,1:0,2:0,3:0,1:1,2:1',parallel)
// pf = exe(f,a1)
func TestPipeline(t *testing.T) {
	serial := func(playable melrose.Sequenceable) melrose.Sequenceable {
		return playable
	}
	parallel := func(playable melrose.Sequenceable) melrose.Sequenceable {
		return playable
	}
	a := Pipeline{
		Target: []interface{}{serial, parallel},
	}
	c := melrose.MustParseChord("C")
	r, err := a.Execute(c)
	if r == nil || err != nil {
		t.Fatal(err)
	}
	t.Log(r.S().Notes)
}

func TestPipelineSerialInterface(t *testing.T) {
	serial := func(playables ...interface{}) interface{} {
		return playables[0]
	}
	a := Pipeline{
		Target: []interface{}{serial},
	}
	c := melrose.MustParseChord("C")
	r, err := a.Execute(c)
	if r == nil || err != nil {
		t.Fatal(err)
	}
	t.Log(r.S().Notes)
}

func TestPipelineParallelInterface(t *testing.T) {
	parallel := func(value interface{}) interface{} {
		return value
	}
	a := Pipeline{
		Target: []interface{}{parallel},
	}
	c := melrose.MustParseChord("C")
	r, err := a.Execute(c)
	if r == nil || err != nil {
		t.Fatal(err)
	}
	t.Log(r.S().Notes)
}

func TestPipelineSerialParallelInterface(t *testing.T) {
	serial := func(playables ...interface{}) interface{} {
		return playables[0]
	}
	parallel := func(value interface{}) interface{} {
		return value
	}
	a := Pipeline{
		Target: []interface{}{serial, parallel},
	}
	c := melrose.MustParseChord("C")
	r, err := a.Execute(c)
	if r == nil || err != nil {
		t.Fatal(err)
	}
	t.Log(r.S().Notes)
}

func TestPipelineRepeat(t *testing.T) {
	repeat := func(howMany int, m interface{}) interface{} {
		return m
	}
	a := Pipeline{
		Target: []interface{}{repeat, 1},
	}
	c := melrose.MustParseChord("C")
	r, err := a.Execute(c)
	if r == nil || err != nil {
		t.Fatal(err)
	}
	t.Log(r.S().Notes)
}

func TestPipelineOctaveMap(t *testing.T) {
	octavemap := func(indices string, m interface{}) interface{} {
		return m
	}
	a := Pipeline{
		Target: []interface{}{octavemap, "1:1"},
	}
	c := melrose.MustParseChord("C")
	r, err := a.Execute(c)
	if r == nil || err != nil {
		t.Fatal(err)
	}
	t.Log(r.S().Notes)
}
