package control

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type SetBPM struct {
	bpm core.HasValue
	ctx core.Context
}

func NewBPM(bpm core.HasValue, ctx core.Context) SetBPM {
	return SetBPM{bpm: bpm, ctx: ctx}
}

// S has the side effect of setting the BPM unless BPM is zero
func (s SetBPM) S() core.Sequence {
	s.Evaluate(s.ctx)
	return core.EmptySequence
}

// Evaluate implements Evaluatable
// performs the set operation
func (s SetBPM) Evaluate(ctx core.Context) error {
	f := core.Float(s.bpm)
	if core.IsDebug() {
		notify.Debugf("control.bpm set %.2f", f)
	}
	if f > 0.0 {
		ctx.Control().SetBPM(float64(f))
	}
	return nil
}

// Inspect implements Inspectable
func (s SetBPM) Inspect(i core.Inspection) {
	i.Properties["bpm"] = fmt.Sprintf("%.2f", core.Float(s.bpm))
}

// Storex implements Storable
func (s SetBPM) Storex() string {
	return fmt.Sprintf("bpm(%v)", s.bpm)
}
