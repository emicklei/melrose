package main

import (
	"reflect"
	"strings"
	"unicode"

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

func completeMe(line string, pos int) (head string, c []string, tail string) {
	start := pos
	// given pos go back to last separator
	runes := []rune(line)
	for i := start; i != 0; i-- {
		if i >= len(runes) {
			continue
		}
		r := runes[i]
		if !unicode.IsLetter(r) && r != '_' {
			start = i + 1
			break
		}
	}
	// vars first
	for k, _ := range varStore.Variables() {
		if strings.HasPrefix(k, line[start:]) {
			c = append(c, k)
		}
	}
	for k, f := range dsl.EvalFunctions(varStore) {
		// TODO start from closest (
		if strings.HasPrefix(k, line[start:]) {
			c = append(c, f.Sample)
		}
	}
	return line[0:start], c, line[pos:]
}
