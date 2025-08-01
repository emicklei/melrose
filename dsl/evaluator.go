package dsl

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/emicklei/melrose/control"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"

	"github.com/expr-lang/expr"
)

type Evaluator struct {
	context core.Context
	funcs   map[string]Function
}

func NewEvaluator(ctx core.Context) *Evaluator {
	return &Evaluator{
		context: ctx,
		funcs:   EvalFunctions(ctx),
	}
}

const fourSpaces = "    "

// Statements are separated by newlines and semicolons.
// If a line is prefixed by one or more TABs then that line is appended to the previous.
// If a line is prefixed by 4 SPACES then that line is appended to the previous.
// Return the result of the last expression or statement.
func (e *Evaluator) EvaluateProgram(source string) (any, error) {
	lines := []string{}
	splitted := strings.Split(source, "\n")

	// indentation
	nrOfLastExpression := -1
	for lineNr, line := range splitted {
		if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, fourSpaces) { // append to previous
			if len(lines) == 0 {
				return nil, errors.New("syntax error, first line cannot start with TAB")
			}
			if nrOfLastExpression+1 != lineNr {
				return nil, fmt.Errorf("syntax error, line with TAB [%d] must be part of expression", lineNr+1)
			}
			lines[len(lines)-1] = withoutTrailingComment(lines[len(lines)-1]) + line // with TAB TODO
			nrOfLastExpression = lineNr
			continue
		}
		lines = append(lines, line)
		nrOfLastExpression = lineNr
	}
	// now, lines dont have leading tabs or 4 spaces

	var lastResult any
	for _, line := range lines {
		// replace all TABs
		line = strings.Replace(line, "\t", " ", -1)
		// whitespaces
		line = strings.TrimSpace(line)
		// empty
		if len(line) == 0 {
			continue
		}
		// comment line
		if strings.HasPrefix(line, "//") {
			continue
		}
		// trailing inline comment
		line = withoutTrailingComment(line)

		for _, statement := range strings.Split(line, ";") {
			// whitespaces
			statement = strings.TrimSpace(statement)
			// empty
			if len(statement) == 0 {
				continue
			}
			result, err := e.evaluateCleanStatement(statement)
			if err != nil {
				return nil, err
			}
			if result != nil {
				lastResult = result
			}
		}
	}
	return lastResult, nil
}

func (e *Evaluator) RecoveringEvaluateProgram(entry string) (any, error) {
	defer func() {
		if err := recover(); err != nil {
			notify.Errorf("%v", err)
			return
		}
	}()
	return e.EvaluateProgram(entry)
}

func (e *Evaluator) evaluateCleanStatement(entry string) (any, error) {
	if value, ok := e.context.Variables().Get(entry); ok {
		return value, nil
	}
	if varName, expression, ok := IsAssignment(entry); ok {
		// variable cannot be named after function
		if _, conflict := e.funcs[varName]; conflict {
			return nil, fmt.Errorf("cannot use variable [%s] because it is a defined function", varName)
		}

		r, err := e.EvaluateExpression(expression)
		if err != nil {
			return nil, err
		}
		return e.handleAssignment(varName, r)
	}

	// evaluate and print
	r, err := e.EvaluateExpression(entry)
	if err != nil {
		return nil, err
	}

	// special case for Loop,Listen,Recording
	if canStop, ok := r.(core.Stoppable); ok {
		varName := e.newSuggestedVariableName(canStop)
		if len(varName) == 0 {
			return nil, fmt.Errorf("this object must assigned to variable name, use e.g. var = %s", canStop.(core.Storable).Storex())
		}
		return e.handleAssignment(varName, r)
	}

	// special case for Evals, put last because Stoppables can be also Evaluatable
	if theEval, ok := r.(core.Evaluatable); ok {
		if err := theEval.Evaluate(e.context); err != nil { // no condition
			return nil, err
		}
	}

	return r, nil
}

// The last expression returned a Stoppable and was not assigned to a variable.
// Generate a name based on the combination of the file and the line (if both given).
func (e *Evaluator) newSuggestedVariableName(stoppable core.Stoppable) string {
	var line int
	if v, ok := e.context.Environment().Load(core.EditorLineEnd); ok {
		line = v.(int)
	} else {
		return ""
	}
	return fmt.Sprintf("%s%d", shortTypeName(stoppable), line)
}

// *core.Loop => loop
func shortTypeName(v any) string {
	if v == nil {
		return "nil"
	}
	parts := strings.Split(fmt.Sprintf("%T", v), ".")
	if len(parts) > 1 {
		return strings.ToLower(parts[len(parts)-1])
	}
	return strings.ToLower(parts[0])
}

func (e *Evaluator) handleAssignment(varName string, r any) (any, error) {
	// check delete
	if r == nil {
		e.context.Variables().Delete(varName)
	} else {
		// special case for Loop
		// if the value is a Loop
		// then if the variable refers to an existing loop
		// 		then change to Target of that loop
		//		else store the loop
		// else store the result
		if theLoop, ok := r.(*core.Loop); ok {
			if storedValue, present := e.context.Variables().Get(varName); present {
				if otherLoop, replaceme := storedValue.(*core.Loop); replaceme {
					otherLoop.SetTarget(theLoop.Target())
					r = otherLoop
				} else {
					// existing variable but not a Loop
					e.context.Variables().Put(varName, theLoop)
				}
			} else {
				// new variable for theLoop
				e.context.Variables().Put(varName, theLoop)
			}
			return r, nil
		}
		// special case for Listen
		// if the value is a Listen
		// then if the variable refers to an existing listen
		// 		then change to Target of that listen
		//		else store the listen
		// else store the result
		if theListen, ok := r.(*control.Listen); ok {
			if storedValue, present := e.context.Variables().Get(varName); present {
				if otherListen, replaceme := storedValue.(*control.Listen); replaceme {
					otherListen.SetTarget(theListen.Target())
					r = otherListen
				} else {
					// existing variable but not a Listen
					e.context.Variables().Put(varName, theListen)
				}
			} else {
				// new variable for theLoop
				e.context.Variables().Put(varName, theListen)
			}
			return r, nil
		}
		// special case for Recording
		// if the value is a Recording
		// then if the variable refers to an existing recording
		// 		then change the Target of that recording
		//		else store the recording
		// else store the result
		if theRecording, ok := r.(*control.Recording); ok {
			if storedValue, present := e.context.Variables().Get(varName); present {
				if storedRecording, replaceme := storedValue.(*control.Recording); replaceme {
					storedRecording.GetTargetFrom(theRecording)
					r = storedRecording
				} else {
					// existing variable but not a Recording
					e.context.Variables().Put(varName, theRecording)
				}
			} else {
				// new variable for theRecording
				e.context.Variables().Put(varName, theRecording)
			}
			return r, nil
		}

		// not a Loop or Listen or Recording
		e.context.Variables().Put(varName, r)
		if aware, ok := r.(core.NameAware); ok {
			aware.VariableName(varName)
		}
	}
	return r, nil
}

// EvaluateExpression returns the result of an expression (entry) using a given store of variables.
// The result is either FunctionResult or a "raw" Go object.
func (e *Evaluator) EvaluateExpression(entry string) (any, error) {
	options := []expr.Option{}
	// since expr 1.14.3
	for _, each := range []string{"join", "repeat", "trim", "replace", "duration", "map"} {
		options = append(options, expr.DisableBuiltin(each))
	}
	env := map[string]any{}
	for k, f := range e.funcs {
		env[k] = f.Func
	}
	for k := range e.context.Variables().Variables() {
		env[k] = variable{Name: k, store: e.context.Variables()}
	}
	options = append(options, expr.Env(env))
	options = append(options, expr.Patch(new(indexedAccessPatcher)))
	addOperatorsTo(env)
	program, err := expr.Compile(entry, append(options, exprOperators()...)...)
	if err != nil {
		// try parsing the entry as a sequence or chord
		// this can be requested from the editor to listen to a part of a sequence,chord,note,progression
		if strings.Contains(entry, "/") {
			if subchord, suberr := core.ParseChord(entry); suberr == nil {
				if notify.IsDebug() {
					notify.Debugf("dsl.evaluate:%s", subchord.Storex())
				}
				return subchord, nil
			}
		}
		// try parsing the entry as a sequence
		if subseq, suberr := core.ParseSequence(entry); suberr == nil {
			if notify.IsDebug() {
				notify.Debugf("dsl.evaluate:%s", subseq.Storex())
			}
			return subseq, nil
		}
		// give up
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

func (e *Evaluator) LookupFunction(fn string) (Function, bool) {
	for name, each := range e.funcs {
		if name == fn {
			return each, true
		}
	}
	return Function{}, false
}

func withoutTrailingComment(s string) string {
	if slashes := strings.Index(s, "//"); slashes != -1 {
		return s[0:slashes]
	}
	return s
}
