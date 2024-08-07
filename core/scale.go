package core

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose/notify"
)

type Scale struct {
	start   Note
	variant int
	// https://en.wikipedia.org/wiki/Scale_(music)#Scales,_steps,_and_intervals
	scaleType string
}

func (s Scale) Storex() string {
	return fmt.Sprintf("scale('%s %s')", s.scaleType, s.start.String())
}

// Replaced is part of Replaceable
func (s Scale) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, s) {
		return to
	}
	return s
}

func NewScale(input string) (Scale, error) {
	s, err := ParseScale(input)
	if err != nil {
		return s, err
	}
	return s, nil
}

func ParseScale(s string) (Scale, error) {
	style := "major" // Ionian mode
	tokens := strings.Split(s, " ")
	if len(tokens) == 2 {
		style = tokens[0]
		s = tokens[1]
	}
	parts := strings.Split(s, "/")
	n, err := ParseNote(parts[0])
	v := Major
	if len(parts) == 2 && parts[1] == "m" {
		v = Minor
	}
	return Scale{start: n, variant: v, scaleType: style}, err
}

var (
	majorScale        = [7]int{0, 2, 4, 5, 7, 9, 11}
	naturalMinorScale = [7]int{0, 1, 3, 5, 7, 8, 10}
	romans            = [7]int{Major, Minor, Minor, Major, Major, Minor, Major}
)

// ChordAt uses one-based index
func (s Scale) ChordAt(index int) Chord {
	if index < 1 || index > 7 {
		notify.Warnf("invalid index for ChordAt, got %d", index)
		return zeroChord()
	}
	if s.variant == Major {
		offset := majorScale[index-1]
		return Chord{start: s.start.Pitched(offset), inversion: Ground, interval: Triad, quality: romans[index-1]}
	}
	// TODO
	return zeroChord()
}

func (s Scale) S() Sequence {
	notes := []Note{}
	steps := majorScale
	if s.variant == Minor {
		steps = naturalMinorScale
	}
	for _, p := range steps {
		notes = append(notes, s.start.Pitched(p))
	}
	return BuildSequence(notes)
}

// one-based
func (s Scale) IndexOfNote(n Note) int {
	steps := majorScale
	if s.variant == Minor {
		steps = naturalMinorScale
	}
	if s.start.Name == n.Name &&
		s.start.Accidental == n.Accidental {
		return 1
	}
	for i := 1; i < len(steps); i++ {
		p := s.start.Pitched(steps[i])
		if p.Name == n.Name &&
			p.Accidental == n.Accidental {
			return i + 1
		}
	}
	return -1 // not found
}

// one-based
func (s Scale) NoteAtIndex(index int) Note {
	steps := majorScale
	if s.variant == Minor {
		steps = naturalMinorScale
	}
	var stepIndex, octaveDelta int
	if index < 1 {
		stepIndex = (index-1)%7 + 7
		octaveDelta = -1 + ((index - 1) / 7)
	} else {
		stepIndex = (index - 1) % 7
		octaveDelta = (index - 1) / 7
	}
	n := s.start.Pitched(steps[stepIndex])
	// n could have a change in the octave
	// so put it back such that delta is correct
	n.Octave = s.start.Octave
	return n.Octaved(octaveDelta)
}
