package api

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

type ServiceImpl struct {
	context   core.Context
	evaluator *dsl.Evaluator
}

func NewService(ctx core.Context) Service {
	return &ServiceImpl{context: ctx, evaluator: dsl.NewEvaluator(ctx)}
}

func (s *ServiceImpl) updateMetadata(file string, lineEnd int, source string) error {
	s.context.Environment().Store(core.WorkingDirectory, filepath.Dir(file))
	// get and store lineEnd end
	breaks := strings.Count(source, "\n")
	if breaks > 0 {
		s.context.Environment().Store(core.EditorLineStart, lineEnd-breaks)
	} else {
		s.context.Environment().Store(core.EditorLineStart, lineEnd)
	}
	return nil
}

func (s *ServiceImpl) CommandInspect(file string, lineEnd int, source string) (interface{}, error) {
	s.updateMetadata(file, lineEnd, source)

	returnValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return nil, err
	}
	// check for function
	if reflect.TypeOf(returnValue).Kind() == reflect.Func {
		if fn, ok := s.evaluator.LookupFunction(source); ok {
			fmt.Fprintf(notify.Console.StandardOut, "%s: %s\n", fn.Title, fn.Description)
		}
	} else {
		core.PrintValue(s.context, returnValue)
	}
	return returnValue, nil
}
func (s *ServiceImpl) CommandPlay(file string, lineEnd int, source string) (interface{}, error) {
	s.updateMetadata(file, lineEnd, source)

	returnValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return nil, err
	}

	if pl, ok := returnValue.(core.Playable); ok {
		notify.Infof("play(%s)", displayString(s.context, pl))
		_ = pl.Play(s.context, time.Now())
	} else {
		// any sequenceable is playable
		if seq, ok := returnValue.(core.Sequenceable); ok {
			notify.Infof("play(%s)", displayString(s.context, seq))
			s.context.Device().Play(
				core.NoCondition,
				seq,
				s.context.Control().BPM(),
				time.Now())
		}
	}

	return returnValue, nil
}
func (s *ServiceImpl) CommandStop(file string, lineEnd int, source string) (interface{}, error) {
	s.updateMetadata(file, lineEnd, source)

	returnValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return nil, err
	}

	if p, ok := returnValue.(core.Stoppable); ok {
		notify.Infof("stop(%s)", displayString(s.context, p))
		return p.Stop(s.context), nil
	}

	return returnValue, nil
}
func (s *ServiceImpl) CommandEvaluate(file string, lineEnd int, source string) (interface{}, error) {
	s.updateMetadata(file, lineEnd, source)

	returnValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return nil, err
	}
	return returnValue, nil
}
func (s *ServiceImpl) CommandKill() error {
	// kill the play and any loop
	dsl.StopAllPlayables(s.context)
	s.context.Device().Reset()
	return nil
}

func (s *ServiceImpl) CommandHover(source string) string {
	// inspect as variable
	value, ok := s.context.Variables().Get(source)
	if ok {
		return core.NewInspect(s.context, value).Markdown()
	}
	// inspect as function
	fun, ok := s.evaluator.LookupFunction(source)
	if ok {
		return fun.Markdown()
	}
	return ""
}

func displayString(ctx core.Context, v interface{}) string {
	name := ctx.Variables().NameFor(v)
	if len(name) == 0 {
		name = core.Storex(v)
	}
	return name
}
