package main

import (
	"log"
	"reflect"
	"strings"

	"github.com/emicklei/melrose"
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
	// if line ends with dot then lookup methods for the value before the dot
	if strings.HasSuffix(line, ".") {
		n := melrose.C()
		for _, each := range availableMethodNames(n, line) {
			c = append(c, line+each)
		}
		log.Println(c)
	} else {
		for k, f := range evalFuncMap {
			// TODO start from closest (
			if strings.HasPrefix(k, strings.ToLower(line)) {
				c = append(c, f.Sample)
			}
		}
	}
	return
}
