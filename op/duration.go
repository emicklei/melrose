package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose"
)

type Duration struct {
	Target    []melrose.Sequenceable
	Parameter float64
}

func NewDuration(checkedParameter float64, target []melrose.Sequenceable) Duration {
	return Duration{
		Target:    target,
		Parameter: checkedParameter,
	}
}

func (d Duration) S() melrose.Sequence {
	target := [][]melrose.Note{}
	source := Join{Target: d.Target}.S().Notes
	for _, eachGroup := range source {
		mappedGroup := []melrose.Note{}
		for _, eachNote := range eachGroup {
			mappedGroup = append(mappedGroup, eachNote.WithDuration(d.Parameter))
		}
		target = append(target, mappedGroup)
	}
	return melrose.Sequence{Notes: target}
}

func (d Duration) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "duration(%f", d.Parameter)
	appendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
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
