package dsl

import "github.com/emicklei/melrose"

type variable struct {
	Name  string
	store *VariableStore
}

func (v variable) Storex() string {
	return v.Name
}

func (v variable) String() string {
	return v.Name
}

func (v variable) S() melrose.Sequence {
	m, ok := v.store.Get(v.Name)
	if !ok {
		return melrose.Sequence{}
	}
	if s, ok := m.(melrose.Sequenceable); ok {
		return s.S()
	}
	return melrose.Sequence{}
}

func (v variable) Value() interface{} {
	m, _ := v.store.Get(v.Name)
	return m
}
