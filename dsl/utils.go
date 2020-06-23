package dsl

import (
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

func StopAllLoops(storage VariableStorage) {
	// stop any running loops
	for k, v := range storage.Variables() {
		if l, ok := v.(*melrose.Loop); ok {
			if l.IsRunning() {
				notify.Print(notify.Infof("stopping loop [%s]", k))
				l.Stop()
			}
		}
	}
}

// Run executes the program (source) and return the value of the last expression or any error while executing.
func Run(ctx *melrose.PlayContext, source string) (interface{}, error) {
	store := NewVariableStore()
	eval := NewEvaluator(store, ctx.LoopControl)

	r, err := eval.EvaluateProgram(source)

	if err != nil {
		return r, err
	}

	// wait until all sounds are played
	for !ctx.AudioDevice.Timeline().IsEmpty() {
		time.Sleep(1 * time.Second)
	}

	return r, err
}
