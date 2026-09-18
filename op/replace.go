package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose/core"
)

// Replace will replace a Sequenceable upon creating a Sequence.
type Replace struct {
	Target   []core.Sequenceable
	From, To core.Sequenceable
}

// S is part of Sequenceable
func (r Replace) S() core.Sequence {
	if len(r.Target) == 0 {
		return core.EmptySequence
	}
	if rep, ok := r.Target[0].(core.Replaceable); ok {
		return rep.Replaced(r.From, r.To).S()
	}
	return r.Target[0].S()
}

// Storex is part of Storable
func (r Replace) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "replace(")
	if len(r.Target) > 0 {
		if st, ok := r.Target[0].(core.Storable); ok {
			fmt.Fprintf(&b, "%s,", st.Storex())
		} else {
			fmt.Fprintf(&b, "%v,", r.Target[0])
		}
	} else {
		fmt.Fprintf(&b, "nil,")
	}
	if st, ok := r.From.(core.Storable); ok {
		fmt.Fprintf(&b, "%s,", st.Storex())
	} else {
		fmt.Fprintf(&b, "%v,", r.From)
	}
	if st, ok := r.To.(core.Storable); ok {
		fmt.Fprintf(&b, "%s)", st.Storex())
	} else {
		fmt.Fprintf(&b, "%v)", r.To)
	}
	return b.String()
}

// Return a new Replace in which any occurrences of "from" are replaced by "to".
func (r Replace) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(r, from) {
		return to
	}
	return Replace{Target: replacedAll(r.Target, from, to), From: r.From, To: r.To}
}
