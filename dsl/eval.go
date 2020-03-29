package dsl

import (
	"regexp"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/emicklei/melrose/notify"
)

type EvaluationResult struct {
	Notification notify.Message
	Result       interface{}
}

func result(r interface{}, m notify.Message) EvaluationResult {
	return EvaluationResult{Notification: m, Result: r}
}

func Evaluate(varStore *VariableStore, entry string) (interface{}, error) {
	// flatten multiline ; expr does not support multiline strings
	entry = strings.Replace(entry, "\n", " ", -1)

	env := map[string]interface{}{}
	for k, f := range EvalFunctions(varStore) { // cache this?
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
