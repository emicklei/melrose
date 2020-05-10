package dsl

import "github.com/antonmedv/expr"

type envMap map[string]interface{}

func (envMap) exprOperators() []expr.Option {
	return []expr.Option{
		expr.Operator("-", "Sub"),
		expr.Operator("+", "Add"),
	}
}
func (envMap) Sub(v variable, i int) int { return i }

func (envMap) Add(l, r interface{}) interface{} {
	if vl, ok := l.(variable); ok {
		return vl.dispatchAdd(r)
	}
	if vr, ok := r.(variable); ok {
		return vr.dispatchAdd(l)
	}
	return nil
}

func (v variable) dispatchAdd(r interface{}) interface{} {
	if vr, ok := r.(variable); ok {
		// int
		il, lok := v.Value().(int)
		ir, rok := vr.Value().(int)
		if lok && rok {
			return il + ir
		}
	}
	if ir, ok := r.(int); ok {
		il, lok := v.Value().(int)
		if lok {
			return il + ir
		}
	}
	return nil
}
