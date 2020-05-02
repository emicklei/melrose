package op

import "github.com/emicklei/melrose"

type Apply struct {
	arguments []interface{}
}

// f = apply(serial,octavemap,'1:-1,3:-1,1:0,2:0,3:0,1:1,2:1',parallel)
// f(a1)

type ApplyFunc func(s melrose.Sequenceable) melrose.Sequenceable

func (a Apply) Func() ApplyFunc {
	return nil
}
