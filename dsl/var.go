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

// At is called from expr after patching []. One-based
func (v variable) At(index int) interface{} {
	m, ok := v.store.Get(v.Name)
	if !ok {
		return nil
	}
	if intArray, ok := m.([]int); ok {
		if index < 1 || index > len(intArray) {
			return nil
		}
		return intArray[index-1]
	}
	if indexable, ok := m.(melrose.Indexable); ok {
		return indexable.At(index)
	}
	if sequenceable, ok := m.(melrose.Sequenceable); ok {
		return melrose.BuildSequence(sequenceable.S().At(index))
	}
	return nil
}

// AtVariable is called from expr after patching [].
func (v variable) AtVariable(index variable) interface{} {
	indexVal := melrose.Int(index)
	if indexVal == 0 {
		return nil
	}
	return v.At(indexVal)
}

type variableArithmetic struct{}

func (variableArithmetic) Sub(v variable, i int) int { return i }
