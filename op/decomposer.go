package op

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emicklei/melrose/core"
)

func DecomposeSequence(s core.Sequence) core.Sequenceable {
	fractions := []int{}
	dynamics := []index2dynamic{}
	notes := []core.Note{}
	for i, eachGroup := range s.Notes {
		if len(eachGroup) == 1 {
			eachNote := eachGroup[0]
			f := int(1.0 / eachNote.Fraction())
			fractions = append(fractions, f)
			s := core.VelocityToDynamic(eachNote.Velocity)
			if s != "" {
				dynamics = append(dynamics, index2dynamic{at: i + 1, dynamic: s}) // 1-based index
			}
			notes = append(notes, eachNote.WithoutDynamic().WithFraction(0.25, false))
		}
	}
	fs := new(bytes.Buffer)
	io.WriteString(fs, "'")
	for i, each := range fractions {
		if i > 0 {
			io.WriteString(fs, " ")
		}
		fmt.Fprintf(fs, "%d", each)
	}
	io.WriteString(fs, "'")
	r := DynamicMap{
		Target: []core.Sequenceable{FractionMap{
			target:   core.BuildSequence(notes),
			fraction: core.On(fs),
		}},
		IndexDynamics: dynamics,
	}
	return r
}
