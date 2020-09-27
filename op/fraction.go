package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Fraction struct {
	Target    []core.Sequenceable
	Parameter float64
}

func NewFraction(checkedParameter float64, target []core.Sequenceable) Fraction {
	return Fraction{
		Target:    target,
		Parameter: checkedParameter,
	}
}

// Return a new Fraction in which any occurrences of "from" are replaced by "to".
func (d Fraction) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if from == core.Sequenceable(d) {
		return to
	}
	return Fraction{Target: replacedAll(d.Target, from, to), Parameter: d.Parameter}
}

func (d Fraction) S() core.Sequence {
	target := [][]core.Note{}
	source := Join{Target: d.Target}.S().Notes
	for _, eachGroup := range source {
		mappedGroup := []core.Note{}
		for _, eachNote := range eachGroup {
			mappedGroup = append(mappedGroup, eachNote.WithFraction(d.Parameter, eachNote.Dotted))
		}
		target = append(target, mappedGroup)
	}
	return core.Sequence{Notes: target}
}

func (d Fraction) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "fraction(%f", d.Parameter)
	appendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (d Fraction) ToNote() core.Note {
	if len(d.Target) == 0 {
		notify.Panic(fmt.Errorf("cannot take note from [%s]", d.Storex()))
	}
	one, ok := d.Target[0].(core.NoteConvertable)
	if !ok {
		notify.Panic(fmt.Errorf("cannot take note from [%v]", one))
	}
	not := one.ToNote()
	return not.WithFraction(d.Parameter, not.Dotted)
}

var validFractionParameterValues = []float64{0.0625, 0.125, 0.25, 0.5, 1, 2, 4, 8, 16}

func CheckFraction(param float64) error {
	match := false
	for _, each := range validFractionParameterValues {
		if each == param {
			match = true
			break
		}
	}
	if !match {
		return fmt.Errorf("fraction parameter [%v] must in %v", param, validFractionParameterValues)
	}
	return nil
}
