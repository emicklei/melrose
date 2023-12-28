package core

import (
	"testing"

	"github.com/expr-lang/expr"
)

func TestIterator_Value(t *testing.T) {
	var l = []interface{}{"C", "D"}
	i := &Iterator{Target: l}
	if got, want := i.Value(), "C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	i.Next()
	if got, want := i.Value(), "D"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := i.Index().Value(), 2; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	i.Next()
	if got, want := i.Value(), "C"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

type loop struct {
	t *testing.T
}

func (l *loop) Index() number {
	l.t.Log("Index call")
	return 1
}

type number int

func TestIndex(t *testing.T) {
	env := map[string]interface{}{
		"a": &loop{t: t},
	}
	_, err := expr.Compile(`a.Index()`, expr.Env(env))
	t.Log(err)
}
