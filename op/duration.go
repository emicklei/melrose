package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Duration struct {
	Target    []core.Sequenceable
	Parameter float64
}

func NewDuration(checkedParameter float64, target []core.Sequenceable) Duration {
	return Duration{
		Target:    target,
		Parameter: checkedParameter,
	}
}

// Return a new Duration in which any occurrences of "from" are replaced by "to".
func (d Duration) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if from == core.Sequenceable(d) {
		return to
	}
	return Duration{Target: replacedAll(d.Target, from, to), Parameter: d.Parameter}
}

func (d Duration) S() core.Sequence {
	target := [][]core.Note{}
	source := Join{Target: d.Target}.S().Notes
	for _, eachGroup := range source {
		mappedGroup := []core.Note{}
		for _, eachNote := range eachGroup {
			mappedGroup = append(mappedGroup, eachNote.WithDuration(d.Parameter))
		}
		target = append(target, mappedGroup)
	}
	return core.Sequence{Notes: target}
}

func (d Duration) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "duration(%f", d.Parameter)
	appendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (d Duration) ToNote() core.Note {
	if len(d.Target) == 0 {
		notify.Panic(fmt.Errorf("cannot take note from [%s]", d.Storex()))
	}
	one, ok := d.Target[0].(core.NoteConvertable)
	if !ok {
		notify.Panic(fmt.Errorf("cannot take note from [%v]", one))
	}
	return one.ToNote().WithDuration(d.Parameter)
}

var validDurationParameterValues = []float64{0.0625, 0.125, 0.25, 0.5, 1, 2, 4, 8, 16}

func CheckDuration(param float64) error {
	match := false
	for _, each := range validDurationParameterValues {
		if each == param {
			match = true
			break
		}
	}
	if !match {
		return fmt.Errorf("duration parameter [%v] must in %v", param, validDurationParameterValues)
	}
	return nil
}
