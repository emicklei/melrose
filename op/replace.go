package op

import (
	"bytes"
	"fmt"
	"github.com/emicklei/melrose/core"
)

// Replace will replace a Sequenceable upon creating a Sequence.
type Replace struct {
	Target   core.Sequenceable
	From, To core.Sequenceable
}

// S is part of Sequenceable
func (r Replace) S() core.Sequence {
	if rep, ok := r.Target.(core.Replaceable); ok {
		return rep.Replaced(r.From, r.To).S()
	}
	return r.Target.S()
}

// Storex is part of Storable
func (r Replace) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "replace(")
	if st, ok := r.Target.(core.Storable); ok {
		fmt.Fprintf(&b, "%s,", st.Storex())
	} else {
		fmt.Fprintf(&b, "%v,", r.Target)
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
