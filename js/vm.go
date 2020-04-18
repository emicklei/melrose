package js

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/dop251/goja"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

// NewVirtualMachineAndStorage returns a Javascript runtime with DSL functions
func NewVirtualMachineAndStorage(store dsl.VariableStorage) (*goja.Runtime, dsl.VariableStorage) {
	vm := goja.New()
	storeWrapper := newAdaptorOn(vm, store)
	for k, v := range dsl.EvalFunctions(storeWrapper) {
		vm.Set(k, v.Func)
	}
	return vm, storeWrapper
}

// LanguageServer can execute DSL statements received over HTTP
type LanguageServer struct {
	vm      *goja.Runtime
	mutex   *sync.Mutex
	store   dsl.VariableStorage
	address string
}

// NewLanguageServer returns a new LanguageService. It is not started.
func NewLanguageServer(vm *goja.Runtime, store dsl.VariableStorage, addr string) *LanguageServer {
	return &LanguageServer{vm: vm, store: store, address: addr, mutex: new(sync.Mutex)}
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
	entry := string(data)
	returnValue, err := l.vm.RunString(entry)
	var response interface{}
	if err != nil {
		// evaluation failed.
		w.WriteHeader(http.StatusInternalServerError)
		notify.Print(notify.Errorf("melrose.js:%s", err.Error()))
		fmt.Println()
		if jserr, ok := err.(*goja.Exception); ok {
			response = errorFrom(jserr)
		} else {
			response = errorFrom(err)
		}
	} else {
		// evaluation was ok.
		if gov := returnValue.Export(); gov != nil {
			// if assignment then make sure our storage knows about the gov.
			if variable, _, ok := dsl.IsAssignment(entry); ok {
				l.store.Put(variable, gov)
			}
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
