package op

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose/core"
)

type DynamicMapper struct {
	Target       core.Sequenceable
	NoteDynamics []string
}

func NewDynamicMapper(s core.Sequenceable, dynamics string) DynamicMapper {
	return DynamicMapper{Target: s, NoteDynamics: strings.Split(dynamics, " ")}
}

func (d DynamicMapper) S() core.Sequence {
	target := [][]core.Note{}
	source := d.Target.S().Notes
	for i, eachGroup := range source {
		if i >= len(d.NoteDynamics) {
			// no change of dynamic
			target = append(target, eachGroup)
		} else {
			mappedGroup := []core.Note{}
			for _, eachNote := range eachGroup {
				mappedGroup = append(mappedGroup, eachNote.WithVelocity(core.ParseVelocity(d.NoteDynamics[i])))
			}
			target = append(target, mappedGroup)
		}
	}
	return core.Sequence{Notes: target}
}

func (d DynamicMapper) Storex() string {
	if s, ok := d.Target.(core.Storable); ok {
		return fmt.Sprintf("dynamicmap(%s,%s)", strings.Join(d.NoteDynamics, " "), s.Storex())
	}
	return ""
}
