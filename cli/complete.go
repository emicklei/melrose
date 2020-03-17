package main

import (
	"reflect"
	"strings"
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
