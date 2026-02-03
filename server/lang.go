package server

import (
	"flag"
	"net/http"

	"github.com/emicklei/melrose/api"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/structexplorer"
)

// LanguageServer can execute DSL statements received over HTTP
type LanguageServer struct {
	context core.Context
	address string
	service api.Service
}

// NewLanguageServer returns a new LanguageService. It is not started.
func NewLanguageServer(ctx core.Context, addr string) *LanguageServer {
	return &LanguageServer{context: ctx, address: addr, service: api.NewService(ctx)}
}

// Start will start a HTTP listener for serving DSL statements
// curl -v -d 'n = note("C")' http://localhost:8118/v1/statements
func (l *LanguageServer) Start() error {
	http.HandleFunc("/v1/statements", l.statementHandler)
	http.HandleFunc("/v1/inspect", l.inspectHandler)
	http.HandleFunc("/v1/notes", l.notesPageHandler)
	http.HandleFunc("/v1/pianoroll", l.pianorollImageHandler)
	http.HandleFunc("/version", l.versionHandler)
	return http.ListenAndServe(l.address, nil)
}

var httpPort = flag.String("http", ":8118", "address on which to listen for HTTP requests")

func Start(ctx core.Context) {
	ls := NewLanguageServer(ctx, *httpPort)
	if len(*httpPort) > 0 {
		// start DSL server
		go ls.Start()
	} else {
		notify.Warnf("empty http flag, skip starting HTTP server")
	}
	if notify.IsDebug() {
		go structexplorer.NewService("ctx", ctx).Start()
	}
}
