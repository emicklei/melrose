package melrose

import (
	"fmt"
	"strings"

	"log"
)

type Scale struct {
	start   Note
	variant int
}

func (s Scale) Storex() string {
	return fmt.Sprintf("scale('%s')", s.start.Storex())
}

func ParseScale(s string) (Scale, error) {
	n, err := ParseNote(s)
	v := Major
	if strings.HasSuffix(s, "/m") {
		v = Minor
	}
	return Scale{start: n, variant: v}, err
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
	for _, p := range steps {
		notes = append(notes, s.start.Pitched(p))
	}
	return BuildSequence(notes)
}
