package op

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose"
)

/**

This is an experiment


**/

type FuncMapper struct {
	Target melrose.Sequenceable
	Func   melrose.MapFunc
}

func (f FuncMapper) S() melrose.Sequenceable {
	return f.Func(f.Target)
}

func (f FuncMapper) Storex() string {
	// if s, ok := r.Target.(Storable); ok {
	// 	return fmt.Sprintf("repeat(%d,%s)", r.Times, s.Storex())
	// }
	return ""
}

// MapFunc returns a MapFunc the creates a new DynamicMapper using the receiver's note dynamics.
func (d DynamicMapper) MapFunc() melrose.MapFunc {
	return func(s melrose.Sequenceable) melrose.Sequenceable {
		return DynamicMapper{
			Target:       s,
			NoteDynamics: d.NoteDynamics,
		}
	}
}

// s1 = sequence("C D")
// s1.dynamic("+ -")
// s1.map(dynamic("+ -"))
// dynamic("+ -",s1)

type DynamicFunction struct {
}

func Dynamic(modifiers string) DynamicFunction {
	return DynamicFunction{}
}

func Repeat(times melrose.Valueable) melrose.MapFunc {
	return func(s melrose.Sequenceable) melrose.Sequenceable {
		return s.S().Repeated(melrose.Int(times))
	}
}

type DynamicOperator struct {
	NoteDynamics []string
}

func NewDynamic(dynamics string) DynamicOperator {
	return DynamicOperator{NoteDynamics: strings.Split(dynamics, " ")}
}

func (d DynamicOperator) Apply(seq melrose.Sequenceable) melrose.Sequenceable {
	target := [][]melrose.Note{}
	source := seq.S().Notes
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

func (d DynamicOperator) Storex(target melrose.Storable) string {
	return fmt.Sprintf("dynamicmap(%s,%s)", strings.Join(d.NoteDynamics, " "), target.Storex())
}
