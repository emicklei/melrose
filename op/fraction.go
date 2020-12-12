package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
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
	f := float32(d.Parameter)
	if f > 1.0 {
		f = 1.0 / f
	}
	target := [][]core.Note{}
	source := Join{Target: d.Target}.S().Notes
	for _, eachGroup := range source {
		mappedGroup := []core.Note{}
		for _, eachNote := range eachGroup {
			mappedGroup = append(mappedGroup, eachNote.WithFraction(f, eachNote.Dotted))
		}
		target = append(target, mappedGroup)
	}
	return core.Sequence{Notes: target}
}

func (d Fraction) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "fraction(%v", d.Parameter)
	core.AppendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (d Fraction) ToNote() (core.Note, error) {
	if len(d.Target) == 0 {
		return core.Rest4, fmt.Errorf("cannot take note from [%s]", d.Storex())
	}
	one, ok := d.Target[0].(core.NoteConvertable)
	if !ok {
		return core.Rest4, fmt.Errorf("cannot take note from [%v]", one)
	}
	not, err := one.ToNote()
	if err != nil {
		return not, err
	}
	return not.WithFraction(1.0/float32(d.Parameter), not.Dotted), nil
}

// TODO used?
var validFractionParameterValues = []float64{1, 2, 4, 8, 16, 0.5, 0.25, 0.125, 0.0625, 0.1666}

// TODO used?
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
