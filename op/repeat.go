package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose"
)

type Repeat struct {
	Target []melrose.Sequenceable
	Times  melrose.Valueable
}

func (r Repeat) S() melrose.Sequence {
	times := melrose.Int(r.Times)
	return Join{Target: r.Target}.S().Repeated(times)
}

func (r Repeat) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "repeat(%v", r.Times)
	appendStorexList(&b, false, r.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}
