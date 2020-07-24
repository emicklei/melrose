package core

import "fmt"

type Iterator struct {
	Index  int
	Target []interface{}
}

func (i *Iterator) Value() interface{} {
	if len(i.Target) == 0 {
		return ""
	}
	if i.Index == len(i.Target) {
		i.Index = 0
	}
	next := i.Target[i.Index]
	i.Index++
	return next
}

func (i Iterator) Storex() string {
	return fmt.Sprintf("iterator(%v)", i.Target)
}
