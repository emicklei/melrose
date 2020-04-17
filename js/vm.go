package js

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	address string
}

// NewLanguageServer returns a new LanguageService. It is not started.
func NewLanguageServer(vm *goja.Runtime, addr string) *LanguageServer {
	return &LanguageServer{vm: vm, address: addr}
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
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	returnValue, err := l.vm.RunString(string(data))
	var response interface{}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		notify.Print(notify.Errorf("JS error:%v", err))
		response = err
	} else {
		response = returnValue.Export()
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(response)
}
