package core

import (
	"bytes"
	"fmt"
	"strings"
)

type ChordSequence struct {
	Chords [][]Chord
}

func MustParseChordSequence(s string) ChordSequence {
	p, err := ParseChordSequence(s)
	if err != nil {
		panic(err)
	}
	return p
}

func ParseChordSequence(input string) (ChordSequence, error) {
	p := ChordSequence{}
	// hack to keep scanning simple, TODO
	splitable := strings.Replace(input, groupOpen, " "+groupOpen+" ", -1)
	splitable = strings.Replace(splitable, groupClose, " "+groupClose+" ", -1)
	parts := strings.Fields(splitable)
	ingroup := false
	var group []Chord
	for _, each := range parts {
		if groupOpen == each {
			ingroup = true
			group = []Chord{}
		} else if groupClose == each {
			ingroup = false
			p.Chords = append(p.Chords, group)
		} else {
			next, err := ParseChord(each)
			if err != nil {
				return p, err
			}
			if ingroup {
				group = append(group, next)
			} else {
				p.Chords = append(p.Chords, []Chord{next})
			}
		}
	}
	return p, nil
}

func (p ChordSequence) S() Sequence {
	notes := [][]Note{}
	for _, eachGroup := range p.Chords {
		if len(eachGroup) == 1 {
			notes = append(notes, eachGroup[0].Notes())
		} else {
			// join all notes of each Chord
			joined := []Note{}
			for _, eachChord := range eachGroup {
				joined = append(joined, eachChord.Notes()...)
			}
			notes = append(notes, joined)
		}
	}
	return Sequence{Notes: notes}
}

// Replaced is part of Replaceable
func (p ChordSequence) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, p) {
		return to
	}
	return p
}

func (p ChordSequence) Storex() string {
	var b bytes.Buffer
	fmt.Fprint(&b, "chordsequence('")
	for i, each := range p.Chords {
		if i > 0 {
			fmt.Fprint(&b, " ")
		}
		if len(each) == 1 {
			fmt.Fprintf(&b, "%s", each[0].String())
		} else {
			fmt.Fprint(&b, "(")
			for j, other := range each {
				if j > 0 {
					fmt.Fprint(&b, " ")
				}
				fmt.Fprintf(&b, "%s", other.String())
			}
			fmt.Fprint(&b, ")")
		}
	}
	fmt.Fprint(&b, "')")
	return b.String()
}

func (p ChordSequence) Inspect(i Inspection) {
	i.Properties[""] = p.S().String()
}
