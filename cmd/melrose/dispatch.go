package main

import (
	"fmt"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

func dispatch(entry string) error {
	if len(entry) == 0 {
		fmt.Println()
		return nil
	}
	if value, ok := varStore.Get(entry); ok {
		fmt.Printf(entry)
		printValue(value)
		return nil
	}
	if variable, expression, ok := dsl.IsAssignment(entry); ok {
		r, err := dsl.Evaluate(varStore, expression)
		if err != nil {
			return err
		}
		if er, ok := r.(dsl.FunctionResult); ok {
			notify.Print(er.Notification)
			// TODO check that we do not use a function name as variable
			varStore.Put(variable, er.Result)
			printValue(er.Result)
		} else {
			// check delete
			if r == nil {
				varStore.Delete(variable)
			} else {
				varStore.Put(variable, r)
				printValue(r)
			}
		}
		return nil
	}
	// evaluate and print
	r, err := dsl.Evaluate(varStore, entry)
	if err != nil {
		return err
	}
	if er, ok := r.(dsl.FunctionResult); ok {
		// info,warn,error
		notify.Print(er.Notification)
		// still can have a result
		printValue(er.Result)
	} else {
		// should not happen
		printValue(r)
	}
	return nil
}

func printValue(v interface{}) {
	if v == nil {
		return
	}
	if s, ok := v.(melrose.Storable); ok {
		fmt.Printf("%s\n", s.Storex())
	} else {
		fmt.Printf("%v\n", v)
	}
}
