package op

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/emicklei/melrose/core"
)

type DynamicMapper struct {
	Target        []core.Sequenceable
	IndexDynamics []index2dynamic
}

type index2dynamic struct {
	at      int // at is one-based
	dynamic string
}

func NewDynamicMapper(slist []core.Sequenceable, dynamics string) DynamicMapper {
	return DynamicMapper{Target: slist, IndexDynamics: parseIndex2Dynamics(dynamics)}
}

func (d DynamicMapper) S() core.Sequence {
	target := [][]core.Note{}
	source := Join{Target: d.Target}.S().Notes
	for _, entry := range d.IndexDynamics {
		if entry.at <= 0 || entry.at > len(source) {
			// invalid offset, skip
			continue
		}
		eachGroup := source[entry.at-1] // at is one-based
		newGroup := []core.Note{}
		for _, eachNote := range eachGroup {
			newGroup = append(newGroup, eachNote.WithVelocity(core.ParseVelocity(entry.dynamic)))
		}
		target = append(target, newGroup)
	}
	return core.Sequence{Notes: target}
}

func (d DynamicMapper) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "dynamic('%s'", formatIndex2Dynamics(d.IndexDynamics))
	appendStorexList(&b, false, d.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// "1:++,2:--"
func parseIndex2Dynamics(s string) (list []index2dynamic) {
	entries := strings.Split(s, ",")
	for _, each := range entries {
		kv := strings.Split(each, ":")
		if len(kv) != 2 { // silent ignore error
			continue
		}
		ik, err := strconv.Atoi(kv[0])
		if err != nil { // silent ignore error TODO
			continue
		}
		list = append(list, index2dynamic{at: ik, dynamic: strings.TrimSpace(kv[1])})
	}
	return
}

// "1:++,2:--"
func formatIndex2Dynamics(list []index2dynamic) string {
	var b bytes.Buffer
	for i, each := range list {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, "%d:%s", each.at, each.dynamic)
	}
	return b.String()
}
