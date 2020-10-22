package core

import (
	"fmt"

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

// S is part of Sequenceable
func (w Watch) S() Sequence {
	beats, bars := w.Context.Control().BeatsAndBars()
	target := fmt.Sprintf("%v", w.Target)
	if v, ok := w.Target.(Valueable); ok {
		target = fmt.Sprintf("%v", v)
	}
	notify.Print(notify.Infof("on bars [%d] beats [%d] called sequence of [%v]", beats, bars, target))
	return EmptySequence
}
