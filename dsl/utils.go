package dsl

import (
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

func StopAllLoops(store *VariableStore) {
	// stop any running loops
	for k, v := range store.Variables() {
		if l, ok := v.(*melrose.Loop); ok {
			notify.Print(notify.Infof("stopping loop [%s]", k))
			l.Stop()
		}
	}
}
