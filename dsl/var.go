package dsl

import (
	"fmt"

	"github.com/emicklei/melrose"
)

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

// dispatchSubFrom  v(l) - r
func (v variable) dispatchSub(r interface{}) interface{} {
	if vr, ok := r.(variable); ok {
		// int
		il, lok := v.Value().(int)
		ir, rok := vr.Value().(int)
		if lok && rok {
			return il - ir
		}
	}
	if ir, ok := r.(int); ok {
		// int
		il, lok := v.Value().(int)
		if lok {
			return il - ir
		}
	}
	return nil
}

// dispatchSubFrom  l - v(r)
func (v variable) dispatchSubFrom(l interface{}) interface{} {
	if vl, ok := l.(variable); ok {
		// int
		il, lok := vl.Value().(int)
		ir, rok := v.Value().(int)
		if lok && rok {
			return il - ir
		}
	}
	if il, ok := l.(int); ok {
		// int
		ir, rok := v.Value().(int)
		if rok {
			return il - ir
		}
	}
	return nil
}

func (v variable) dispatchAdd(r interface{}) interface{} {
	if vr, ok := r.(variable); ok {
		// int
		il, lok := v.Value().(int)
		ir, rok := vr.Value().(int)
		if lok && rok {
			return il + ir
		}
	}
	if ir, ok := r.(int); ok {
		il, lok := v.Value().(int)
		if lok {
			return il + ir
		}
	}
	return nil
}
