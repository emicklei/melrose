package dsl

import (
	"strings"
	"testing"

	"github.com/emicklei/melrose"
)

func newTestEvaluator() *Evaluator {
	return NewEvaluator(testContext())
}

func checkError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func eval(t *testing.T, expression string) interface{} {
	t.Helper()
	lp := new(melrose.TestLooper)
	lp.SetBIAB(4)
	ctx := melrose.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     lp,
	}
	r, err := NewEvaluator(ctx).EvaluateExpression(expression)
	checkError(t, err)
	return r
}

func mustError(t *testing.T, expression string, substring string) {
	t.Helper()
	lp := new(melrose.TestLooper)
	lp.SetBIAB(4)
	ctx := melrose.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     lp,
	}
	_, err := NewEvaluator(ctx).EvaluateExpression(expression)
	if err == nil {
		t.Fatal("error expected")
	}
	if !strings.Contains(err.Error(), substring) {
		t.Fatalf("error message should contain [%s] but it [%s]", substring, err.Error())
	}
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
