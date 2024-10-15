package core

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/notify"
)

type Print struct {
	Context Context
	Target  interface{}
}

func (w Print) Play(ctx Context, at time.Time) error {
	w.S()
	return nil
}

func (w Print) Evaluate(ctx Context) error {
	// TODO check c?
	w.S()
	return nil
}

// S is part of Sequenceable
func (w Print) S() Sequence {
	beats, bars := w.Context.Control().BeatsAndBars()
	in := NewInspect(w.Context, "", w.Target)
	if bars > 0 {
		in.Properties["bar"] = bars
	}
	if beats > 0 {
		in.Properties["beat"] = beats
	}
	notify.Infof("%s", in.String())
	return EmptySequence
}

// Storex is part of Storable
func (w Print) Storex() string {
	return fmt.Sprintf("print(%s)", Storex(w.Target))
}
