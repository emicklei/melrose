package api

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl"
	midifile "github.com/emicklei/melrose/midi/file"
	"github.com/emicklei/melrose/notify"
	"github.com/expr-lang/expr/file"
)

type ServiceImpl struct {
	context   core.Context
	evaluator *dsl.Evaluator
}

func NewService(ctx core.Context) Service {
	return &ServiceImpl{context: ctx, evaluator: dsl.NewEvaluator(ctx)}
}

func (s *ServiceImpl) Context() core.Context { return s.context }

func (s *ServiceImpl) ChangeDefaultDeviceAndChannel(isInput bool, deviceID int, channel int) error {
	// TODO handle channel
	var err error
	if isInput {
		err = s.context.Device().HandleSetting("midi.in", []interface{}{deviceID})
	} else {
		err = s.context.Device().HandleSetting("midi.out", []interface{}{deviceID})
	}
	if err != nil {
		notify.Errorf("change device/channel failed:%s", err.Error())
	}
	return nil
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
	s.context.Environment().Store(core.EditorLineEnd, lineEnd)
	return nil
}

func (s *ServiceImpl) CommandInspect(file string, lineEnd int, source string) (interface{}, error) {
	s.updateMetadata(file, lineEnd, source)

	lastValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return nil, patchFilelocation(err, lineEnd)
	}
	if lastValue == nil {
		core.PrintValue(s.context, nil)
		return lastValue, nil
	}
	// check for function
	if reflect.TypeOf(lastValue).Kind() == reflect.Func {
		if fn, ok := s.evaluator.LookupFunction(source); ok {
			fmt.Fprintf(notify.Console.StandardOut, "%s: %s\n", fn.Title, fn.Description)
		}
	} else {
		core.PrintValue(s.context, lastValue)
	}
	return lastValue, nil
}

type CommandPlayResponse struct {
	EndTime          time.Time
	ExpressionResult any
}

func (s *ServiceImpl) CommandPlay(file string, lineEnd int, source string) (CommandPlayResponse, error) {
	s.updateMetadata(file, lineEnd, source)

	programResult, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return CommandPlayResponse{}, patchFilelocation(err, lineEnd)
	}
	var endTime time.Time
	if pl, ok := programResult.(core.Playable); ok {
		notify.Infof("play(%s)", displayString(s.context, pl))
		endTime = pl.Play(s.context, time.Now())
	} else {
		// unvalue if needed
		if u, ok := programResult.(core.HasValue); ok {
			programResult = u.Value()
		}
		// any sequenceable is playable
		if seq, ok := programResult.(core.Sequenceable); ok {
			notify.Infof("play(%s)", displayString(s.context, seq))
			endTime = s.context.Device().Play(
				core.NoCondition,
				seq,
				s.context.Control().BPM(),
				time.Now())
		} else {
			notify.Debugf("not sequenceable:%v [%T]", programResult, programResult)
		}
	}
	return CommandPlayResponse{EndTime: endTime, ExpressionResult: programResult}, nil
}
func (s *ServiceImpl) CommandStop(file string, lineEnd int, source string) (interface{}, error) {
	s.updateMetadata(file, lineEnd, source)

	returnValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return nil, patchFilelocation(err, lineEnd)
	}

	if p, ok := returnValue.(core.Stoppable); ok {
		notify.Infof("stopping(%s)", displayString(s.context, p))
		return returnValue, p.Stop(s.context)
	}

	return returnValue, nil
}
func (s *ServiceImpl) CommandEvaluate(file string, lineEnd int, source string) (interface{}, error) {
	s.updateMetadata(file, lineEnd, source)

	returnValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return nil, patchFilelocation(err, lineEnd)
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
		return core.NewInspect(s.context, source, value).Markdown()
	}
	// inspect as function
	fun, ok := s.evaluator.LookupFunction(source)
	if ok {
		return fun.Markdown()
	}
	return ""
}

func (s *ServiceImpl) CommandMIDISample(filename string, lineEnd int, source string) ([]byte, error) {
	s.updateMetadata(filename, lineEnd, source)

	returnValue, err := s.evaluator.EvaluateProgram(source)
	if err != nil {
		return []byte{}, patchFilelocation(err, lineEnd)
	}
	buffer := new(bytes.Buffer)
	err = midifile.ExportOn(buffer, returnValue, s.context.Control().BPM(), s.context.Control().BIAB())
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func (s *ServiceImpl) ListDevices() []core.DeviceDescriptor {
	return s.context.Device().ListDevices()
}

func displayString(ctx core.Context, v interface{}) string {
	name := ctx.Variables().NameFor(v)
	if len(name) == 0 {
		name = core.Storex(v)
	}
	return name
}

func patchFilelocation(err error, lineEnd int) error {
	// patch Location of error
	if fe, ok := err.(*file.Error); ok {
		fe.Line = fe.Line - 1 + lineEnd
		return fe
	}
	return err
}
