package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Dynamic struct {
	Target   []core.Sequenceable
	Emphasis core.Valueable
}

// Replaced returns a new Dynamic in which any occurrences of "from" are replaced by "to".
func (d Dynamic) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if from == core.Sequenceable(d) {
		return to
	}
	return Dynamic{Target: replacedAll(d.Target, from, to), Emphasis: d.Emphasis}
}

func (d Dynamic) S() core.Sequence {
	target := [][]core.Note{}
	source := Join{Target: d.Target}.S().Notes
	for _, eachGroup := range source {
		mappedGroup := []core.Note{}
		for _, eachNote := range eachGroup {
			// emphasis is a string or int
			e := d.Emphasis.Value()
			n := eachNote
			if s, ok := e.(string); ok {
				if core.ParseVelocity(s) == -1 {
					notify.Warnf("[op.Dynamic] invalid velocity:%s", s)
					return core.EmptySequence
				}
				n = eachNote.WithDynamic(s)
			} else if v, ok := e.(int); ok {
				n = eachNote.WithVelocity(v)
			}
			mappedGroup = append(mappedGroup, n)
		}
		target = append(target, mappedGroup)
	}
	return core.Sequence{Notes: target}
}

func CheckDynamic(emphasis string) error {
	if core.ParseVelocity(emphasis) == -1 {
		return fmt.Errorf("[op.Dynamic] dynamic parameter [%v] must in %v", emphasis, "{+,++,+++,++++,-,--,---,----,o}")
	}
	return nil
}

func (d Dynamic) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "dynamic(%s", core.Storex(d.Emphasis))
	core.AppendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (d Dynamic) ToNote() (core.Note, error) {
	if len(d.Target) == 0 {
		return core.Rest4, fmt.Errorf("[op.Dynamic] cannot take note from [%s]", d.Storex())
	}
	one, ok := d.Target[0].(core.NoteConvertable)
	if !ok {
		return core.Rest4, fmt.Errorf("[op.Dynamic] cannot take note from [%v]", one)
	}
	not, err := one.ToNote()
	if err != nil {
		return not, err
	}
	// emphasis is a string or int
	e := d.Emphasis.Value()
	if s, ok := e.(string); ok {
		if core.ParseVelocity(s) == -1 {
			return not, fmt.Errorf("[op.Dynamic] invalid velocity:%s", s)
		}
		not = not.WithDynamic(s)
	} else if v, ok := e.(int); ok {
		not = not.WithVelocity(v)
	}
	return not, nil
}
