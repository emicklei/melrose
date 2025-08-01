package core

import "fmt"

type Iterator struct {
	index  int
	name   string
	Target []any
}

// VariableName is part of NameAware
func (i *Iterator) VariableName(n string) {
	i.name = n
}

// Value is part of HasValue
func (i *Iterator) Value() any {
	if len(i.Target) == 0 {
		return nil
	}
	return i.Target[i.index]
}

// Index returns the current index of the iterator as a HasValue ; 1-based
func (i *Iterator) Index() HasValue {
	return ValueFunction{
		StoreString: fmt.Sprintf("%s.Index()", i.name),
		Function: func() any {
			return i.getindex() + 1
		}}
}

func (i *Iterator) getindex() int { return i.index }

// S is part of Sequenceable
func (i *Iterator) S() Sequence {
	if len(i.Target) == 0 {
		return EmptySequence
	}
	v := i.Target[i.index]
	if s, ok := v.(Sequenceable); ok {
		return s.S()
	}
	return EmptySequence
}

// Next is part of Nextable
func (i *Iterator) Next() any { // TODO return value needed?
	if len(i.Target) == 0 {
		return nil
	}
	if i.index+1 == len(i.Target) {
		i.index = 0
	} else {
		i.index++
	}
	return i.Value()
}

// Storex is part of Storable
func (i Iterator) Storex() string {
	return fmt.Sprintf("iterator(%v)", i.Target)
}

// Inspect is part of Inspectable
func (i Iterator) Inspect(in Inspection) {
	in.Properties["index"] = i.index + 1
	in.Properties["value"] = i.Value()
	in.Properties["length"] = len(i.Target)
}
