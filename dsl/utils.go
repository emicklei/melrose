package dsl

import (
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

func StopAllLoops(storage VariableStorage) {
	// stop any running loops
	for k, v := range storage.Variables() {
		if l, ok := v.(*melrose.Loop); ok {
			notify.Print(notify.Infof("stopping loop [%s]", k))
			l.Stop()
		}
	}
}
