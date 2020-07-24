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
	Target  Sequenceable
}

// S is part of Sequenceable
func (w Watch) S() Sequence {
	beats, bars := w.Context.Control().BeatsAndBars()
	target := fmt.Sprintf("%v", w.Target)
	st, ok := w.Target.(Storable)
	if ok {
		target = st.Storex()
	}
	notify.Print(notify.Infof("on bars [%d] beats [%d] called sequence of [%s]", beats, bars, target))
	return w.Target.S()
}
