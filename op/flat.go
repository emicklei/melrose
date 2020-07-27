package op

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type Flattener struct {
	Target core.Sequenceable
}

func (f Flattener) S() core.Sequence {
	return f.Target.S()
}

func (f Flattener) Storex() string {
	if st, ok := f.Target.(core.Storable); ok {
		return fmt.Sprintf("flat(%s)", st.Storex())
	}
	return "?"
}
