package op

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose"
)

type DynamicMapper struct {
	Target       melrose.Sequenceable
	NoteDynamics []string
}

func NewDynamicMapper(s melrose.Sequenceable, dynamics string) DynamicMapper {
	return DynamicMapper{Target: s, NoteDynamics: strings.Split(dynamics, " ")}
}

func (d DynamicMapper) S() melrose.Sequence {
	target := [][]melrose.Note{}
	source := d.Target.S().Notes
	for i, eachGroup := range source {
		if i >= len(d.NoteDynamics) {
			// no change of dynamic
			target = append(target, eachGroup)
		} else {
			mappedGroup := []melrose.Note{}
			for _, eachNote := range eachGroup {
				mappedGroup = append(mappedGroup, eachNote.ModifiedVelocity(melrose.ParseVelocity(d.NoteDynamics[i])))
			}
			target = append(target, mappedGroup)
		}
	}
	return melrose.Sequence{Notes: target}
}

func (d DynamicMapper) Storex() string {
	if s, ok := d.Target.(melrose.Storable); ok {
		return fmt.Sprintf("dynamicmap(%s,%s)", strings.Join(d.NoteDynamics, " "), s.Storex())
	}
	return ""
}
