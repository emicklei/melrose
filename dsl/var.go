package dsl

import "github.com/emicklei/melrose"

type VariableStorage interface {
	NameFor(value interface{}) string
	Get(key string) (interface{}, bool)
	Put(key string, value interface{})
	Delete(key string)
	Variables() map[string]interface{}
}

type variable struct {
	Name  string
	store VariableStorage
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

type variableArithmetic struct{}

func (variableArithmetic) Sub(v variable, i int) int { return i }
