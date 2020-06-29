package dsl

import (
	"fmt"
	"github.com/emicklei/melrose/core"
)

type variable struct {
	Name  string
	store core.VariableStorage
}

func (v variable) Storex() string {
	return v.Name
}

func (v variable) String() string {
	return fmt.Sprintf("var %s", v.Name)
}

func (v variable) S() core.Sequence {
	m, ok := v.store.Get(v.Name)
	if !ok {
		return core.EmptySequence
	}
	if s, ok := m.(core.Sequenceable); ok {
		return s.S()
	}
	return core.EmptySequence
}

// Replaced is part of Replaceable
func (v variable) Replaced(from, to core.Sequenceable) core.Sequenceable {
	if core.IsIdenticalTo(from, v) {
		return to
	}
	currentValue := v.Value()
	if currentS, ok := currentValue.(core.Sequenceable); ok {
		if core.IsIdenticalTo(from, currentS) {
			return to
		}
	}
	if rep, ok := currentValue.(core.Replaceable); ok {
		return rep.Replaced(from, to)
	}
	return v
}

func (v variable) Value() interface{} {
	m, _ := v.store.Get(v.Name)
	return m
}
