package core

import (
	"fmt"

	"github.com/emicklei/melrose/notify"
)

type ChordProgression struct {
	root     Valueable
	sequence Valueable
}

func NewChordProgression(root, sequence Valueable) ChordProgression {
	return ChordProgression{root: root, sequence: sequence}
}

// Replaced is part of Replaceable
func (c ChordProgression) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, c) {
		return to
	}
	return c
}

// Storex is part of Storable
func (c ChordProgression) Storex() string {
	return fmt.Sprintf("progression(%s,%s)", Storex(c.root), Storex(c.sequence))
}

// S is part of Sequenceable
func (c ChordProgression) S() Sequence {
	cs, ok := c.root.Value().(string)
	if !ok {
		notify.Warnf("chord progression root must be string, type: %T", c.root.Value())
		return EmptySequence
	}
	sc, err := ParseScale(cs)
	if err != nil {
		notify.Warnf("chord progression root must use scale notation, error: %v", err)
		return EmptySequence
	}
	input, ok := c.sequence.Value().(string)
	if !ok {
		notify.Warnf("chord progression must be string, type: %T", c.sequence.Value())
		return EmptySequence
	}
	p := newFormatParser(input)
	chords, err := p.parseChordProgression(sc)
	if err != nil {
		notify.Warnf("parsing progression failed, error: %v", err)
		return EmptySequence
	}
	j := EmptySequence
	for _, each := range chords {
		j = j.SequenceJoin(each.S())
	}
	return j
}
