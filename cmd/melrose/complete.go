package main

import (
	"reflect"
	"strings"

	"github.com/emicklei/melrose/dsl"
)

func availableMethodNames(v interface{}, prefix string) (list []string) {
	rt := reflect.TypeOf(v)
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		if strings.HasPrefix(m.Name, prefix) {
			list = append(list, m.Name)
		}
	}
	return
}

func completeMe(line string) (c []string) {
	// vars first
	for k, _ := range varStore.Variables() {
		if strings.HasPrefix(k, line) {
			c = append(c, k)
		}
	}
	for k, f := range dsl.EvalFunctions(varStore) {
		// TODO start from closest (
		if strings.HasPrefix(k, line) {
			c = append(c, f.Sample)
		}
	}
	return
}
