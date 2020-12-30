package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func (l *LanguageServer) inspectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		notify.Console.Warnf("HTTP method not allowed:%s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	debug := r.URL.Query().Get("debug") == "true" || core.IsDebug()
	if debug {
		notify.Debugf("service.http: %s", r.URL.String())
	}
	// get token
	defer r.Body.Close()
	type tokenHolder struct {
		Token string
	}
	th := new(tokenHolder)
	if err := json.NewDecoder(r.Body).Decode(th); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if debug {
		notify.Debugf("http.request.body %#v", th)
	}
	type markdownHolder struct {
		MarkdownString string
	}
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	msg := l.markdownOnInspecting(th.Token)
	err := enc.Encode(markdownHolder{MarkdownString: msg})
	if debug {
		notify.Debugf("service.http.response.MarkdownString: %s", msg)
	}
	if err != nil {
		notify.Console.Errorf("inspect failed:%v\n", err)
	}
}

func (l *LanguageServer) markdownOnInspecting(token string) string {
	// inspect as variable
	value, ok := l.context.Variables().Get(token)
	if ok {
		return core.NewInspect(l.context, value).Markdown()
	}
	fun, ok := l.evaluator.LookupFunction(token)
	if ok {
		return fmt.Sprintf("%s\n\n%s", fun.Description, fun.Template)
	}
	return ""
}
