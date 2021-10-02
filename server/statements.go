package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func (l *LanguageServer) statementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		notify.Console.Warnf("HTTP method not allowed:%s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	file := query.Get("file")
	debug := query.Get("debug") == "true" || core.IsDebug()
	if debug {
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
	// get expression source
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	source := string(data)

	if debug {
		notify.Debugf("http.request.body %s", source)
	}
	defer r.Body.Close()

	var evalResult interface{}
	action := query.Get("action")
	switch action {
	case "kill":
		evalResult = l.service.CommandKill()
	case "inspect":
		if ret, err := l.service.CommandInspect(file, line, source); err != nil {
			evalResult = err
		} else {
			evalResult = ret
		}
	case "play":
		if ret, err := l.service.CommandPlay(file, line, source); err != nil {
			evalResult = err
		} else {
			evalResult = ret
		}
	case "stop":
		if ret, err := l.service.CommandStop(file, line, source); err != nil {
			evalResult = err
		} else {
			evalResult = ret
		}
	case "eval":
		if ret, err := l.service.CommandEvaluate(file, line, source); err != nil {
			evalResult = err
		} else {
			evalResult = ret
		}
	default:
		evalResult = fmt.Errorf("unknown command:%s", query.Get("action"))
	}
	if _, ok := evalResult.(error); ok {
		// evaluation failed.
		w.WriteHeader(http.StatusBadRequest)
	}
	response := resultFrom(file, line, evalResult)
	w.Header().Set("content-type", "application/json")
	w.Header().Set("access-control-allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(response)
	if err != nil {
		notify.NewErrorf("error:%v\n", err)
		return
	}
	if response.IsError {
		notify.Print(notify.NewError(response.Object.(error)))
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
	Type         string      `json:"type"`
	IsError      bool        `json:"is-error"`
	IsStoppeable bool        `json:"stoppable"`
	Message      string      `json:"message"`
	Filename     string      `json:"file"`
	Line         int         `json:"line"`
	Column       int         `json:"column"`
	Object       interface{} `json:"object"`
}

func resultFrom(filename string, line int, val interface{}) evaluationResult {
	t := fmt.Sprintf("%T", val)
	_, isStoppable := val.(core.Stoppable)
	if err, ok := val.(error); ok {
		return evaluationResult{
			Type:         t,
			IsError:      true,
			IsStoppeable: isStoppable,
			Filename:     filename,
			Message:      err.Error(),
			Line:         line,
			Object:       val,
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
	return evaluationResult{
		Type:         t,
		IsError:      false,
		IsStoppeable: isStoppable,
		Filename:     filename,
		Line:         line,
		Message:      msg}
}
