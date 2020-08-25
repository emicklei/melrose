package dsl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/emicklei/melrose/core"
)

// ImportProgram runs a script from a file
func ImportProgram(ctx core.Context, filename string) error {
	pwd, ok := ctx.Environment()[core.WorkingDirectory]
	if !ok {
		pwd = ""
	}
	fullName := filepath.Join(pwd, filename)
	data, err := ioutil.ReadFile(fullName)
	if err != nil {
		abs, _ := filepath.Abs(fullName)
		return fmt.Errorf("unable to read file[%s] :%v", abs, err)
	}
	eval := NewEvaluator(ctx)
	_, err = eval.EvaluateProgram(string(data))
	return err
}
