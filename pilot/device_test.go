package pilot

import (
	"testing"

	"github.com/emicklei/melrose"
)

func TestPlaySequence(t *testing.T) {
	p, err := Open()
	if err != nil {
		t.Log(err)
		return
	}
	defer p.Close()
	p.SetBeatsPerMinute(140)
	s := melrose.MustParseSequence("C D E")
	p.Play(s, true)
	s2 := melrose.MustParseSequence("(C D E)")
	p.Play(s2, true)
}
