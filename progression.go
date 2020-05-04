package melrose

import (
	"bytes"
	"fmt"
	"strings"
)

type Progression struct {
	Chords [][]Chord
}

func MustParseProgression(s string) Progression {
	p, err := ParseProgression(s)
	if err != nil {
		panic(err)
	}
	return p
}

func ParseProgression(input string) (Progression, error) {
	p := Progression{}
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

func (p Progression) S() Sequence {
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

func (p Progression) Storex() string {
	var b bytes.Buffer
	fmt.Fprint(&b, "progression('")
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
