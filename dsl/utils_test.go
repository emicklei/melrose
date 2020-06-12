package dsl

import (
	"testing"

	"github.com/emicklei/melrose"
)

func newTestEvaluator() *Evaluator {
	v := NewVariableStore()
	e := NewEvaluator(v, melrose.NoLooper)
	return e
}

func checkError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func eval(t *testing.T, expression string) interface{} {
	melrose.Context().LoopControl = new(melrose.TestLooper)
	r, err := newTestEvaluator().EvaluateExpression(expression)
	checkError(t, err)
	return r
}

func checkStorex(t *testing.T, r interface{}, storex string) {
	t.Helper()
	if s, ok := r.(melrose.Storable); ok {
		if got, want := s.Storex(), storex; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	} else {
		t.Errorf("result is not Storable : [%v:%T]", r, r)
	}
}
