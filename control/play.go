package control

import (
	"bytes"
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
)

// Play represents play() and sync()
type Play struct {
	ctx    core.Context
	target []core.Sequenceable
	sync   bool
}

func NewPlay(ctx core.Context, list []core.Sequenceable, playInSync bool) Play {
	return Play{
		ctx:    ctx,
		target: list,
		sync:   playInSync,
	}
}

// Play is part of core.Playable
func (p Play) Play(ctx core.Context, at time.Time) error {
	cond := core.NoCondition
	if with, ok := ctx.(core.Conditional); ok {
		cond = with.Condition()
	}
	for _, each := range p.target {
		end := p.ctx.Device().Play(cond, each, p.ctx.Control().BPM(), at)
		if !p.sync {
			// play after each other
			at = end
		}
	}
	return nil
}

// Evaluate implements Evaluatable
// performs the set operation
func (p Play) Evaluate(ctx core.Context) error {
	moment := time.Now()
	cond := core.NoCondition
	if with, ok := ctx.(core.Conditional); ok {
		cond = with.Condition()
	}
	for _, each := range p.target {
		end := p.ctx.Device().Play(cond, each, p.ctx.Control().BPM(), moment)
		if !p.sync {
			// play after each other
			moment = end
		}
	}
	return nil
}

// Storex implements Storable
func (p Play) Storex() string {
	var b bytes.Buffer
	if p.sync {
		fmt.Fprintf(&b, "sync(")
	} else {
		fmt.Fprintf(&b, "play(")
	}
	core.AppendStorexList(&b, true, p.target)
	fmt.Fprintf(&b, ")")
	return b.String()
}
