package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose"
)

type Replace struct {
	Target   melrose.Sequenceable
	From, To melrose.Sequenceable
}

func (r Replace) S() melrose.Sequence {
	if rep, ok := r.Target.(melrose.Replaceable); ok {
		fmt.Println("replace inside target")
		return rep.Replaced(r.From, r.To).S()
	}
	fmt.Printf("not replaced inside target (%T)", r.Target)
	return r.Target.S()
}

// Storex is part of Storable
func (r Replace) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "replace(")
	if st, ok := r.Target.(melrose.Storable); ok {
		fmt.Fprintf(&b, "%s,", st.Storex())
	} else {
		fmt.Fprintf(&b, "%v,", r.Target)
	}
	if st, ok := r.From.(melrose.Storable); ok {
		fmt.Fprintf(&b, "%s,", st.Storex())
	} else {
		fmt.Fprintf(&b, "%v,", r.From)
	}
	if st, ok := r.To.(melrose.Storable); ok {
		fmt.Fprintf(&b, "%s)", st.Storex())
	} else {
		fmt.Fprintf(&b, "%v)", r.To)
	}
	return b.String()
}
