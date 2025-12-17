package op

import (
	"fmt"
	"reflect"

	"github.com/emicklei/melrose/core"
)

// IfCondition represents a conditional operation that selects between two sequences
// In the DSL, it is represented as when(condition, then, else). "if" is deprecated.
type IfCondition struct {
	Condition core.HasValue
	Then      core.Sequenceable
	Else      core.Sequenceable
}

func (i IfCondition) S() core.Sequence {
	b, ok := i.Condition.Value().(bool)
	if !ok {
		return i.Else.S()
	}
	if !b {
		return i.Else.S()
	}
	return i.Then.S()
}

func (i IfCondition) Storex() string {
	if reflect.DeepEqual(i.Else, core.EmptySequence) {
		return fmt.Sprintf("if(%s,%s)", core.Storex(i.Condition), core.Storex(i.Then))
	}
	return fmt.Sprintf("if(%s,%s,%s)", core.Storex(i.Condition), core.Storex(i.Then), core.Storex(i.Else))
}

// Replaced is part of Replaceable
func (i IfCondition) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(i, from) {
		return to
	}
	if core.IsIdenticalTo(i.Then, from) {
		return IfCondition{Condition: i.Condition, Then: to, Else: i.Else}
	}
	if core.IsIdenticalTo(i.Else, from) {
		return IfCondition{Condition: i.Condition, Then: i.Then, Else: to}
	}
	repThen := i.Then
	if r, ok := repThen.(core.Replaceable); ok {
		repThen = r.Replaced(from, to)
	}
	repElse := i.Else
	if r, ok := repElse.(core.Replaceable); ok {
		repElse = r.Replaced(from, to)
	}
	return IfCondition{Condition: i.Condition, Then: repThen, Else: repElse}
}
