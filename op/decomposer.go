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
	groups := [][]core.Note{}
	for i, eachGroup := range s.Notes {
		if len(eachGroup) == 1 {
			eachNote := eachGroup[0]
			f := int(1.0 / eachNote.Fraction())
			fractions = append(fractions, f)
			s := core.VelocityToDynamic(eachNote.Velocity)
			if s != "" {
				dynamics = append(dynamics, index2dynamic{at: i + 1, dynamic: s}) // 1-based index
			}
			note := eachNote.WithoutDynamic().WithFraction(0.25, false)
			groups = append(groups, []core.Note{note})
		} else {
			f := int(1.0 / eachGroup[0].Fraction()) // first will decide
			fractions = append(fractions, f)
			s := core.VelocityToDynamic(eachGroup[0].Velocity) // first will decide
			if s != "" {
				dynamics = append(dynamics, index2dynamic{at: i + 1, dynamic: s}) // 1-based index
			}
			group := groupWithoutDynamic(eachGroup)
			group = groupWithFraction(group, 0.25, false)
			groups = append(groups, group)
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
			target:   core.Sequence{Notes: groups},
			fraction: core.On(fs),
		}},
		IndexDynamics: dynamics,
	}
	return r
}

func groupWithoutDynamic(g []core.Note) (list []core.Note) {
	for _, each := range g {
		list = append(list, each.WithoutDynamic())
	}
	return
}
func groupWithFraction(g []core.Note, fraction float32, dotted bool) (list []core.Note) {
	for _, each := range g {
		list = append(list, each.WithFraction(fraction, dotted))
	}
	return
}
