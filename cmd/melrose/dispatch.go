package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/emicklei/melrose"
)

func dispatch(entry string) error {
	if len(entry) == 0 {
		fmt.Println()
		return nil
	}
	if value, ok := varStore.Get(entry); ok {
		printValue(value)
		return nil
	}
	if variable, expression, ok := isAssignment(entry); ok {
		r, err := eval(expression)
		if err != nil {
			return err
		}
		// TODO check that we do not use a function name as variable
		varStore.Put(variable, r)
		printValue(r)
		return nil
	}
	// evaluate and print
	r, err := eval(entry)
	if err != nil {
		return err
	}
	// if reflect.TypeOf(r).I
	// 	printWarning(fmt.Sprintf("did you mean %s()?", entry))
	// 	return nil
	// }
	printValue(r)
	return nil
}

func printValue(v interface{}) {
	if v == nil {
		fmt.Println()
		return
	}
	if s, ok := v.(melrose.Storable); ok {
		fmt.Printf("%s\n", s.Storex())
	} else {
		fmt.Printf("%v\n", v)
	}
}

func eval(entry string) (interface{}, error) {
	// flatten multiline ; expr does not support multiline strings
	entry = strings.Replace(entry, "\n", " ", -1)

	env := map[string]interface{}{}
	for k, f := range evalFunctions() {
		env[k] = f.Func
	}
	for k, v := range varStore.Variables() {
		env[k] = v
	}
	program, err := expr.Compile(entry, expr.Env(env))
	if err != nil {
		return nil, err
	}
	return expr.Run(program, env)
}

// https://regex101.com/
var assignmentRegex = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.*)$`)

// [ ]a[]=[]note('c')
func isAssignment(entry string) (varname string, expression string, ok bool) {
	sanitized := strings.TrimSpace(entry)
	res := assignmentRegex.FindAllStringSubmatch(sanitized, -1)
	if len(res) != 1 {
		return "", "", false
	}
	if len(res[0]) != 3 {
		return "", "", false
	}
	return res[0][1], res[0][2], true
}
