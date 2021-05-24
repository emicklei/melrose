package dsl

import (
	"fmt"

	"github.com/emicklei/melrose/core"
)

type variable struct {
	Name  string
	store core.VariableStorage
}

// Storex is part of Storable
func (v variable) Storex() string {
	return v.Name
}

func (v variable) String() string {
	currentValue, _ := v.store.Get(v.Name)
	return fmt.Sprintf("var %s = %v", v.Name, currentValue)
}

func (v variable) Inspect(i core.Inspection) {
	i.Properties["name"] = v.Name
	currentValue, ok := v.store.Get(v.Name)
	if !ok {
		i.Properties["error"] = "missing value"
		return
	}
	i.Properties["val"] = core.Storex(currentValue)
	if insp, ok := currentValue.(core.Inspectable); ok {
		insp.Inspect(i)
	}
}

// S is part of Sequenceable
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

func (v variable) Index() interface{} {
	m, _ := v.store.Get(v.Name)
	if h, ok := m.(core.HasIndex); ok {
		return h.Index()
	}
	return 0 // signals invalid index, public indexables are 1-based
}
