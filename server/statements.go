package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/antonmedv/expr/file"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func (l *LanguageServer) statementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	l.context.Environment()[core.WorkingDirectory] = filepath.Dir(query.Get("file"))

	debug := query.Get("debug") == "true"
	if debug && !core.IsDebug() {
		core.ToggleDebug()
		defer core.ToggleDebug()
	}
	if core.IsDebug() {
		notify.Debugf("service.http: %s", r.URL.String())
	}
	// get line
	line := 1
	lineString := query.Get("line")
	if len(lineString) > 0 {
		if i, err := strconv.Atoi(lineString); err == nil {
			line = i
		}
	}
	// get expression
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if query.Get("action") == "kill" {
		// kill the play and any loop
		l.context.Control().Reset()
		l.context.Device().Reset()
		return
	}
	returnValue, err := l.evaluator.EvaluateProgram(string(data))
	var response evaluationResult
	if err != nil {
		// evaluation failed.
		w.WriteHeader(http.StatusInternalServerError)
		response = resultFrom(line, err)
	} else {
		// evaluation was ok.

		if query.Get("action") == "inspect" {
			// check for function
			if reflect.TypeOf(returnValue).Kind() == reflect.Func {
				if fn, ok := l.evaluator.LookupFunction(string(data)); ok {
					fmt.Fprintf(notify.Console.StandardOut, "%s: %s\n", fn.Title, fn.Description)
				}
			} else {
				core.PrintValue(l.context, returnValue)
			}
		}

		// check if play was requested and is playable
		if query.Get("action") == "play" {
			// first check Playable
			if pl, ok := returnValue.(core.Playable); ok {
				_ = pl.Play(l.context)
			} else {
				// any sequenceable is playable
				if s, ok := returnValue.(core.Sequenceable); ok {
					l.context.Device().Play(
						s,
						l.context.Control().BPM(),
						time.Now())
				}
			}
		}
		// loop operation
		if query.Get("action") == "begin" {
			if lp, ok := returnValue.(*core.Loop); ok {
				if !lp.IsRunning() {
					l.context.Control().StartLoop(lp)
				}
			}
			// ignore if not Loop
		}
		// loop operation
		if query.Get("action") == "end" {
			if lp, ok := returnValue.(*core.Loop); ok {
				if lp.IsRunning() {
					lp.Stop()
				}
			}
			// ignore if not Loop
		}

		// setter TODO make interface?
		if set, ok := returnValue.(core.SetBPM); ok {
			set.S()
		}

		response = resultFrom(line, returnValue)
	}
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(response)
	if err != nil {
		notify.Errorf("error:%v\n", err)
		return
	}
	if response.IsError {
		notify.Print(notify.Error(response.Object.(error)))
	} else {
		core.PrintValue(l.context, response.Object)
	}
	if debug {
		// doit again
		buf := bytes.Buffer{}
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "\t")
		err = enc.Encode(response)
		notify.Debugf("http.response: %s error=%v", buf.String(), err)
	}
}

type evaluationResult struct {
	Type     string      `json:"type"`
	IsError  bool        `json:"is-error"`
	Message  string      `json:"message"`
	Filename string      `json:"file"`
	Line     int         `json:"line"`
	Column   int         `json:"column"`
	Object   interface{} `json:"object"`
}

func resultFrom(line int, val interface{}) evaluationResult {
	t := fmt.Sprintf("%T", val)
	if err, ok := val.(error); ok {
		// patch Location of error
		if fe, ok := err.(*file.Error); ok {
			fe.Location.Line = fe.Location.Line - 1 + line
		}
		return evaluationResult{
			Type:     t,
			IsError:  true,
			Filename: "yours.mel",
			Message:  err.Error(),
			Line:     line,
			Object:   val,
		}
	}
	// no error
	var msg string
	if stor, ok := val.(core.Storable); ok {
		msg = stor.Storex()
	} else {
		msg = fmt.Sprintf("%v", val)
	}
	// no Object if ok
	return evaluationResult{Type: t, IsError: false, Message: msg}
}
