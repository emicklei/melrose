package main

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

var assignmentRegex := regex.MustCompile(`^[a-z]+\[[0-9]+\]=.*$`)

func dispatch(entry string) error {
	if value, ok := memory[entry]; ok {
		fmt.Println(value)
		return nil
	}
	// is assignment?
	// TODO not correct
	if strings.Contains(entry, "=") {
		parts := strings.Split(entry, "=")
		variable := strings.TrimSpace(parts[0])
		expression := parts[1]
		r, err := eval(expression)
		if err != nil {
			return err
		}
		memory[variable] = r
		fmt.Println(r)
		return nil
	}
	// evaluate and print
	r, err := eval(entry)
	if err != nil {
		return err
	}
	fmt.Println(r)
	return nil
}

func eval(entry string) (interface{}, error) {
	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf,
	}
	for k, v := range memory {
		env[k] = v
	}
	program, err := expr.Compile(entry, expr.Env(env))
	if err != nil {
		return nil, err
	}
	return expr.Run(program, env)
}
