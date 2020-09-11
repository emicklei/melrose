package op

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Dynamic struct {
	Target   []core.Sequenceable
	Emphasis string
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
			mappedGroup = append(mappedGroup, eachNote.WithDynamic(d.Emphasis))
		}
		target = append(target, mappedGroup)
	}
	return core.Sequence{Notes: target}
}

var validDynamics = regexp.MustCompile("[+{0,3}]|[-{0,3}]|[0?]")

func CheckDynamic(emphasis string) error {
	if !validDynamics.Match([]byte(emphasis)) {
		return fmt.Errorf("dynamic parameter [%v] must in %v", emphasis, "{+,++,+++,-,--,---,0}")
	}
	return nil
}

func (d Dynamic) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "dynamic('%s'", d.Emphasis)
	appendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (d Dynamic) ToNote() core.Note {
	if len(d.Target) == 0 {
		notify.Panic(fmt.Errorf("cannot take note from [%s]", d.Storex()))
	}
	one, ok := d.Target[0].(core.NoteConvertable)
	if !ok {
		notify.Panic(fmt.Errorf("cannot take note from [%v]", one))
	}
	not := one.ToNote()
	return not.WithDynamic(d.Emphasis)
}
