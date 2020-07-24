package core

import "fmt"

type Iterator struct {
	Index  int
	Target []interface{}
}

//  Value is part of Valueable
func (i *Iterator) Value() interface{} {
	if len(i.Target) == 0 {
		return ""
	}
	return i.Target[i.Index]
}

// Next is part of Nextable
func (i *Iterator) Next() interface{} { // TODO return value needed?
	if len(i.Target) == 0 {
		return nil
	}
	if i.Index+1 == len(i.Target) {
		i.Index = 0
	} else {
		i.Index++
	}
	return i.Value()
}

// Storex is part of Storable
func (i Iterator) Storex() string {
	return fmt.Sprintf("iterator(%v)", i.Target)
}

// Inspect is part of Inspectable
func (i Iterator) Inspect(in Inspection) {
	in.Properties["index"] = i.Index + 1
	in.Properties["length"] = len(i.Target)
}
