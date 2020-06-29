package dsl

import (
	"github.com/emicklei/melrose/core"
	"time"

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

// Run executes the program (source) and returns the value of the last expression or any error while executing.
func Run(device core.AudioDevice, source string) (interface{}, error) {
	ctx := &core.PlayContext{
		AudioDevice:     device,
		VariableStorage: NewVariableStore(),
	}
	ctx.LoopControl = core.NewBeatmaster(ctx, 120)
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
