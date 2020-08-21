package server

import (
	"net/http"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/dsl"
)

// LanguageServer can execute DSL statements received over HTTP
type LanguageServer struct {
	context   core.Context
	address   string
	evaluator *dsl.Evaluator
}

// NewLanguageServer returns a new LanguageService. It is not started.
func NewLanguageServer(ctx core.Context, addr string) *LanguageServer {
	return &LanguageServer{context: ctx, address: addr, evaluator: dsl.NewEvaluator(ctx)}
}

// Start will start a HTTP listener for serving DSL statements
// curl -v -d 'n = note("C")' http://localhost:8118/v1/statements
func (l *LanguageServer) Start() error {
	http.HandleFunc("/v1/statements", l.statementHandler)
	http.HandleFunc("/v1/inspect", l.inspectHandler)
	return http.ListenAndServe(l.address, nil)
}
