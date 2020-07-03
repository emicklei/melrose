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
	Target  Sequenceable
}

func (w Watch) S() Sequence {
	beats, bars := w.Context.Control().BeatsAndBars()
	notify.Print(notify.Infof("on bars [%d] beats [%d] called sequence of [%v]", beats, bars, w.Target))
	return w.Target.S()
}
