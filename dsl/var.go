package dsl

import (
	"fmt"

	"github.com/emicklei/melrose"
)

type variable struct {
	Name  string
	store melrose.VariableStorage
}

func (v variable) Storex() string {
	return v.Name
}

func (v variable) String() string {
	return fmt.Sprintf("var %s", v.Name)
}

func (v variable) S() melrose.Sequence {
	m, ok := v.store.Get(v.Name)
	if !ok {
		return melrose.EmptySequence
	}
	if s, ok := m.(melrose.Sequenceable); ok {
		return s.S()
	}
	return melrose.EmptySequence
}

// Replaced is part of Replaceable
func (v variable) Replaced(from, to melrose.Sequenceable) melrose.Sequenceable {
	if melrose.IsIdenticalTo(from, v) {
		return to
	}
	currentValue := v.Value()
	if currentS, ok := currentValue.(melrose.Sequenceable); ok {
		if melrose.IsIdenticalTo(from, currentS) {
			return to
		}
	}
	if rep, ok := currentValue.(melrose.Replaceable); ok {
		return rep.Replaced(from, to)
	}
	return v
}

func (v variable) Value() interface{} {
	m, _ := v.store.Get(v.Name)
	return m
}
