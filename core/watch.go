package core

import (
	"github.com/emicklei/melrose/notify"
)

var debugEnabled = false

func IsDebug() bool {
	return debugEnabled
}

func ToggleDebug() bool {
	debugEnabled = !debugEnabled
	return debugEnabled
}

type Watch struct {
	Context Context
	Target  interface{}
}

func (w Watch) Evaluate(c Condition) error {
	// TODO check c?
	w.S()
	return nil
}

// S is part of Sequenceable
func (w Watch) S() Sequence {
	beats, bars := w.Context.Control().BeatsAndBars()
	in := NewInspect(w.Context, w.Target)
	in.Properties["bar"] = bars
	in.Properties["beat"] = beats
	notify.Print(notify.Infof("%s", in.String()))
	return EmptySequence
}
