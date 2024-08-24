package dsl

import (
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/mpg"
)

// TEMP

func (v variable) Rate(rate any) variable {
	val, ok := v.store.Get(v.Name)
	if !ok {
		return v
	}
	e, ok := val.(*mpg.Euclidean)
	if !ok {
		return v
	}
	e.Rate = core.On(rate)
	return v
}
