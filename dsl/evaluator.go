package dsl

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/emicklei/melrose"
)

type Evaluator struct {
	store       VariableStorage
	funcs       map[string]Function
	loopControl melrose.LoopController
}

func NewEvaluator(store VariableStorage, loopControl melrose.LoopController) *Evaluator {
	return &Evaluator{
		store:       store,
		funcs:       EvalFunctions(store, loopControl),
		loopControl: loopControl,
	}
}

func (e *Evaluator) EvaluateProgram(source string) error {
	// TODO
	return errors.New("not implemented")
}

func (e *Evaluator) EvaluateStatement(entry string) (interface{}, error) {
	// flatten multiline ; expr does not support multiline strings
	entry = strings.Replace(entry, "\n", " ", -1)
	entry = strings.Replace(entry, "\t", " ", -1)
	if slashes := strings.Index(entry, "//"); slashes != -1 {
		entry = entry[0:slashes]
	}

	if len(entry) == 0 {
		fmt.Println()
		return nil, nil
	}
	// check comment line
	if strings.HasPrefix(entry, "//") {
		return nil, nil
	}
	if value, ok := e.store.Get(entry); ok {
		return value, nil
	}
	if variable, expression, ok := IsAssignment(entry); ok {
		r, err := e.EvaluateExpression(expression)
		if err != nil {
			return nil, err
		}
		// check delete
		if r == nil {
			e.store.Delete(variable)
		} else {
			// special case for Loop
			// if the value is a Loop
			// then if the variable refers to an existing loop
			// 		then change to Target of that loop
			//		else store the loop and run it
			// else store the result
			if theLoop, ok := r.(*melrose.Loop); ok {
				if storedValue, present := e.store.Get(variable); present {
					if otherLoop, replaceme := storedValue.(*melrose.Loop); replaceme {
						otherLoop.Target = theLoop.Target
						e.loopControl.Begin(otherLoop)
						// if !otherLoop.IsRunning() {
						// 	otherLoop.Start(melrose.CurrentDevice())
						// }
					} else {
						// existing variable but not a Loop
						e.store.Put(variable, theLoop)
					}
				} else {
					// new variable for theLoop
					e.store.Put(variable, theLoop)
					e.loopControl.Begin(theLoop)
					//theLoop.Start(melrose.CurrentDevice())
				}
			} else {
				// not a Loop
				e.store.Put(variable, r)
			}
		}
		return r, nil
	}
	// evaluate and print
	r, err := e.EvaluateExpression(entry)
	// special case for Loop
	if theLoop, ok := r.(*melrose.Loop); ok {
		return nil, fmt.Errorf("cannot run an unidentified Loop, use myLoop = %s", theLoop.Storex())
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}

// EvaluateExpression returns the result of an expression (entry) using a given store of variables.
// The result is either FunctionResult or a "raw" Go object.
func (e *Evaluator) EvaluateExpression(entry string) (interface{}, error) {
	env := map[string]interface{}{}
	for k, f := range e.funcs {
		env[k] = f.Func
	}
	for k, _ := range e.store.Variables() {
		env[k] = variable{Name: k, store: e.store}
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
