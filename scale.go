package melrose

import (
	"fmt"
	"strings"

	"log"
)

type Scale struct {
	start   Note
	variant int
	octaves int
}

func (s Scale) Storex() string {
	return fmt.Sprintf("scale(%d,'%s')", s.octaves, s.start.String())
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
	n, err := ParseNote(s)
	v := Major
	if strings.HasSuffix(s, "/m") {
		v = Minor
	}
	return Scale{start: n, variant: v, octaves: 1}, err
}

func MustParseScale(s string) Scale {
	sc, err := ParseScale(s)
	if err != nil {
		log.Fatal(err)
	}
	return sc
}

var (
	majorScale        = [7]int{0, 2, 4, 5, 7, 9, 11}
	naturalMinorScale = [7]int{0, 1, 3, 5, 7, 8, 10}
)

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
	return s.start
}
