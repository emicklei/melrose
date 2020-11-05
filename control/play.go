package control

import (
	"bytes"
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/op"
)

type Play struct {
	ctx    core.Context
	target []core.Sequenceable
}

func NewPlay(ctx core.Context, list []core.Sequenceable) Play {
	return Play{
		ctx:    ctx,
		target: list,
	}
}

// Evaluate implements Evaluatable
// performs the set operation
func (p Play) Evaluate() error {
	moment := time.Now()
	for _, each := range p.target {
		moment = p.ctx.Device().Play(each, p.ctx.Control().BPM(), moment)
	}
	return nil
}

// Storex implements Storable
func (p Play) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "play(")
	op.AppendStorexList(&b, true, p.target)
	fmt.Fprintf(&b, ")")
	return b.String()
}
