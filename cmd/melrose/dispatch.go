package main

import (
	"fmt"
	"strings"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
)

func dispatch(entry string) error {
	if len(entry) == 0 {
		fmt.Println()
		return nil
	}
	// check comment line
	if strings.HasPrefix(entry, "//") {
		return nil
	}
	if value, ok := varStore.Get(entry); ok {
		printValue(value)
		return nil
	}
	if variable, expression, ok := dsl.IsAssignment(entry); ok {
		r, err := dsl.Evaluate(varStore, expression)
		if err != nil {
			return err
		}
		// check delete
		if r == nil {
			varStore.Delete(variable)
		} else {
			varStore.Put(variable, r)
		}
		return nil
	}
	// evaluate and print
	r, err := dsl.Evaluate(varStore, entry)
	if err != nil {
		return err
	}
	printValue(r)
	return nil
}

func printValue(v interface{}) {
	if v == nil {
		return
	}
	if s, ok := v.(melrose.Storable); ok {
		fmt.Printf("\033[94m(%T)\033[0m %s\n", v, s.Storex())
	} else {
		fmt.Printf("\033[94m(%T)\033[0m %v\n", v, v)
	}
}
