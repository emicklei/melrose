package server

import (
	"flag"
	"net/http"

	"github.com/emicklei/melrose/api"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/jonasfj/go-localtunnel"
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
	http.HandleFunc("/version", l.versionHandler)
	return http.ListenAndServe(l.address, nil)
}

var localTunnelSubdomain = flag.String("tunnel", "", "use this as a subdomain to route all commands")

var httpPort = flag.String("http", ":8118", "address on which to listen for HTTP requests")

func Start(ctx core.Context) {
	ls := NewLanguageServer(ctx, *httpPort)
	if len(*httpPort) > 0 {
		// start DSL server
		go ls.Start()
	} else {
		notify.Warnf("empty http flag, skip starting HTTP server")
	}
	if len(*localTunnelSubdomain) > 0 {
		// start local tunnel
		listener, err := localtunnel.Listen(localtunnel.Options{Subdomain: *localTunnelSubdomain})
		if err != nil {
			notify.Errorf("unable to create localtunnel listener:%v", err)
			return
		}
		mux := new(http.ServeMux)
		mux.HandleFunc("/v1/statements", ls.statementHandler)
		mux.HandleFunc("/v1/inspect", ls.inspectHandler)
		mux.HandleFunc("/version", ls.versionHandler)
		server := http.Server{Handler: mux}
		// Handle request from localtunnel
		go server.Serve(listener)
		notify.Infof("opened local tunnel with public address: https://%s.loca.lt", *localTunnelSubdomain)
	}
}
