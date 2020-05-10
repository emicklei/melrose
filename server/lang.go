package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

// LanguageServer can execute DSL statements received over HTTP
type LanguageServer struct {
	store     dsl.VariableStorage
	address   string
	evaluator *dsl.Evaluator
}

// NewLanguageServer returns a new LanguageService. It is not started.
func NewLanguageServer(store dsl.VariableStorage, loopControl melrose.LoopController, addr string) *LanguageServer {
	return &LanguageServer{store: store, address: addr, evaluator: dsl.NewEvaluator(store, loopControl)}
}

// Start will start a HTTP listener for serving DSL statements
// curl -v -d 'n = note("C")' http://localhost:8118/v1/statements
func (l *LanguageServer) Start() error {
	http.HandleFunc("/v1/statements", l.statementHandler)
	return http.ListenAndServe(l.address, nil)
}

func (l *LanguageServer) statementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	trace := r.FormValue("trace") == "true"
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	returnValue, err := l.evaluator.EvaluateProgram(string(data))
	var response interface{}
	if err != nil {
		// evaluation failed.
		w.WriteHeader(http.StatusInternalServerError)
		notify.Print(notify.Errorf("yourfile.mel: %s\n", err.Error()))
		response = errorFrom(err)
	} else {
		// evaluation was ok.

		// check if play was requested and is playable
		if r.FormValue("action") == "play" {
			if s, ok := returnValue.(melrose.Sequenceable); ok {
				melrose.Context().AudioDevice.Play(
					s,
					melrose.Context().LoopControl.BPM(),
					time.Now())
			}
			if lp, ok := returnValue.(*melrose.Loop); ok {
				v := l.store.NameFor(lp)
				notify.Print(notify.Infof("begin loop: %v", v))
				melrose.Context().LoopControl.Begin(lp)
			}
		}
		if r.FormValue("action") == "begin" {
			if lp, ok := returnValue.(*melrose.Loop); ok {
				v := l.store.NameFor(lp)
				notify.Print(notify.Infof("begin loop: %v", v))
				melrose.Context().LoopControl.Begin(lp)
			}
			// ignore if not Loop
		}
		if r.FormValue("action") == "end" {
			if lp, ok := returnValue.(*melrose.Loop); ok {
				v := l.store.NameFor(lp)
				notify.Print(notify.Infof("end loop: %v", v))
				lp.Stop()
			}
			// ignore if not Loop
		}
		if r.FormValue("action") == "kill" {
			// kill the play and any loop
			melrose.Context().LoopControl.Reset()
			melrose.Context().AudioDevice.Reset()
		}
		response = resultFrom(returnValue)
	}
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(response)
	if trace {
		log.Printf("[melrose.trace] %#v\n", response)
	}
}

type evaluationResult struct {
	Type   string `json:"type"`
	Object interface{}
}

func resultFrom(val interface{}) evaluationResult {
	return evaluationResult{Type: fmt.Sprintf("%T", val), Object: val}
}

type evaluationError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}

func errorFrom(err error) evaluationError {
	return evaluationError{Type: fmt.Sprintf("%T", err), Message: err.Error()}
}
