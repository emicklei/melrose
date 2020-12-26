package dsl

import (
	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
)

func StopAllPlayables(context core.Context) {
	// stop any running playables
	for k, v := range context.Variables().Variables() {
		if l, ok := v.(core.Playable); ok {
			notify.Print(notify.Infof("stopping: %s = %s", k, core.Storex(l)))
			_ = l.Stop(context)
		}
	}
}
