package dsl

import (
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

func StopAllLoops(context melrose.Context) {
	// stop any running loops
	for k, v := range context.Variables().Variables() {
		if l, ok := v.(*melrose.Loop); ok {
			if l.IsRunning() {
				notify.Print(notify.Infof("stopping loop [%s]", k))
				l.Stop()
			}
		}
	}
}

// Run executes the program (source) and return the value of the last expression or any error while executing.
func Run(ctx melrose.Context, source string) (interface{}, error) {
	eval := NewEvaluator(ctx)

	r, err := eval.EvaluateProgram(source)

	if err != nil {
		return r, err
	}

	// wait until all sounds are played
	for !ctx.Device().Timeline().IsEmpty() {
		time.Sleep(1 * time.Second)
	}

	return r, err
}
