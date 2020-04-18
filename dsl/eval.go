package dsl

import (
	"regexp"
	"strings"

	"github.com/antonmedv/expr"
)

// Evaluate returns the result of an expression (entry) using a given store of variables.
// The result is either FunctionResult or a "raw" Go object.
func Evaluate(storage VariableStorage, entry string) (interface{}, error) {
	// flatten multiline ; expr does not support multiline strings
	entry = strings.Replace(entry, "\n", " ", -1)

	env := map[string]interface{}{}
	for k, f := range EvalFunctions(storage) { // cache this?
		env[k] = f.Func
	}
	for k, _ := range storage.Variables() {
		env[k] = variable{Name: k, store: storage}
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
func IsAssignment(entry string) (varname string, expression string, ok bool) {
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
