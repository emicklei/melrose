package js

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/dop251/goja"
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

// NewVirtualMachine returns a Javascript runtime with DSL functions
func NewVirtualMachine() *goja.Runtime {
	vm := goja.New()

	// TODO can we use the DSL functions directly?

	vm.Set("seq", func(s string) melrose.Sequence {
		return melrose.MustParseSequence(s)
	})
	vm.Set("play", func(s melrose.Sequenceable) interface{} {
		melrose.CurrentDevice().Play(s, true)
		return nil
	})
	vm.Set("loop", func(s ...melrose.Sequenceable) *melrose.Loop {
		l := &melrose.Loop{Target: melrose.Join{List: s}}
		return l
	})
	vm.Set("run", func(l *melrose.Loop) *melrose.Loop {
		l.Start(melrose.CurrentDevice())
		return l
	})
	vm.Set("note", func(s string) melrose.Note {
		return melrose.MustParseNote(s)
	})
	return vm
}

// LanguageServer can execute DSL statements received over HTTP
type LanguageServer struct {
	vm      *goja.Runtime
	mutex   *sync.Mutex
	address string
}

// NewLanguageServer returns a new LanguageService. It is not started.
func NewLanguageServer(vm *goja.Runtime, addr string) *LanguageServer {
	return &LanguageServer{vm: vm, address: addr, mutex: new(sync.Mutex)}
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
	// the vm is for single thread usage
	l.mutex.Lock()
	defer l.mutex.Unlock()
	returnValue, err := l.vm.RunString(string(data))
	var response interface{}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		notify.Print(notify.Errorf("melrose.js:%s", err.Error()))
		if jserr, ok := err.(*goja.Exception); ok {
			response = errorFrom(jserr)
		} else {
			response = errorFrom(err)
		}
	} else {
		if gov := returnValue.Export(); gov != nil {
			response = resultFrom(gov)
		}
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

// ReferenceError: bpm is not defined at <eval>:1:4(3)
func errorFrom(err error) evaluationError {
	return evaluationError{Type: fmt.Sprintf("%T", err), Message: err.Error()}
}
