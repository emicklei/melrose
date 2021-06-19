package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Fraction struct {
	Target    []core.Sequenceable
	Parameter core.HasValue
}

func NewFraction(parameter core.HasValue, target []core.Sequenceable) Fraction {
	return Fraction{
		Target:    target,
		Parameter: parameter,
	}
}

// Return a new Fraction in which any occurrences of "from" are replaced by "to".
func (d Fraction) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if from == core.Sequenceable(d) {
		return to
	}
	return Fraction{Target: replacedAll(d.Target, from, to), Parameter: d.Parameter}
}

func (d Fraction) floatParameter() float32 {
	f := core.Float(d.Parameter)
	if f > 1.0 {
		f = 1.0 / f
	}
	return f
}

func (d Fraction) S() core.Sequence {
	f := d.floatParameter()
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
	fmt.Fprintf(&b, "fraction(%s", core.Storex(d.Parameter))
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
	return not.WithFraction(d.floatParameter(), not.Dotted), nil
}
