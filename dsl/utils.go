package dsl

import (
	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
)

func StopAllPlayables(context core.Context) {
	// stop any running playables
	for k, v := range context.Variables().Variables() {
		if s, ok := v.(core.Stoppable); ok {
			notify.Infof("stopping: %s = %s", k, core.Storex(s))
			_ = s.Stop(context)
		}
	}
}
