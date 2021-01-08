package core

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose/notify"
)

type Scale struct {
	start   Note
	variant int
	octaves int
}

func (s Scale) Storex() string {
	return fmt.Sprintf("scale(%d,'%s')", s.octaves, s.start.String())
}

// Replaced is part of Replaceable
func (s Scale) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(from, s) {
		return to
	}
	return s
}

func NewScale(octaves int, input string) (Scale, error) {
	s, err := ParseScale(input)
	if err != nil {
		return s, err
	}
	s.octaves = octaves
	return s, nil
}

func ParseScale(s string) (Scale, error) {
	parts := strings.Split(s, "/")
	n, err := ParseNote(parts[0])
	v := Major
	if len(parts) == 2 && parts[1] == "m" {
		v = Minor
	}
	return Scale{start: n, variant: v, octaves: 1}, err
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
	for o := 0; o < s.octaves; o++ {
		for _, p := range steps {
			notes = append(notes, s.start.Pitched(p+(o*12)))
		}
	}
	return BuildSequence(notes)
}

func (s Scale) At(index int) interface{} {
	// TODO
	return s.start
}
