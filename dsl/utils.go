package dsl

import (
	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
)

func StopAllLoops(context core.Context) {
	// stop any running loops
	for k, v := range context.Variables().Variables() {
		if l, ok := v.(*core.Loop); ok {
			if l.IsRunning() {
				notify.Print(notify.Infof("stopping loop [%s]", k))
				l.Stop()
			}
		}
	}
}
