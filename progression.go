package melrose

import (
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
