package op

import (
	"github.com/emicklei/melrose"
)

type WideChord struct {
	Target melrose.Valueable
}

func (w WideChord) Storex() string { return "nil" }

func (w WideChord) S() melrose.Sequence {
	return melrose.Sequence{Notes: [][]melrose.Note{w.Notes()}}
}

func (w WideChord) Notes() (wide []melrose.Note) {
	c, ok := w.Target.Value().(melrose.Chord)
	if !ok {
		return
	}
	notes := c.Notes()
	// method 1
	wide = append(wide, notes[len(notes)-1].Octaved(-1))
	wide = append(wide, notes...)
	wide = append(wide, notes[0].Octaved(1))

	// method 2

	return
}
