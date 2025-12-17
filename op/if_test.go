package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestIf_S(t *testing.T) {
	c := core.On(true)
	th := core.MustParseSequence("C")
	el := core.MustParseSequence("D")
	i := IfCondition{Condition: c, Then: th, Else: el}
	if got, want := i.S().Storex(), "sequence('C')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	i.Condition = core.On(false)
	if got, want := i.S().Storex(), "sequence('D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	i.Condition = core.On("not a boolean")
	if got, want := i.S().Storex(), "sequence('D')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestIf_Storex(t *testing.T) {
	c := core.On(true)
	th := core.MustParseSequence("C")
	el := core.MustParseSequence("D")
	i := IfCondition{Condition: c, Then: th, Else: el}
	if got, want := i.Storex(), "if(true,sequence('C'),sequence('D'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	i.Else = core.EmptySequence
	if got, want := i.Storex(), "if(true,sequence('C'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestIf_Replaced(t *testing.T) {
	c := core.On(true)
	th := core.MustParseSequence("C")
	el := core.MustParseSequence("D")
	i := IfCondition{Condition: c, Then: th, Else: el}

	if core.IsIdenticalTo(i, th) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(i.Replaced(th, el).(IfCondition).Then, el) {
		t.Error("then not replaced")
	}
	if !core.IsIdenticalTo(i.Replaced(el, th).(IfCondition).Else, th) {
		t.Error("else not replaced")
	}
	if !core.IsIdenticalTo(i.Replaced(i, th), th) {
		t.Error("if not replaced")
	}
	rep := i.Replaced(core.EmptySequence, th)
	if !core.IsIdenticalTo(i, rep) {
		t.Error("should be same")
	}
	i.Then = i // self
	if !core.IsIdenticalTo(i.Replaced(i, th), th) {
		t.Error("if not replaced")
	}
}
