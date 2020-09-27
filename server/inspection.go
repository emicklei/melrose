package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/emicklei/melrose/core"
)

func (l *LanguageServer) inspectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	type markdownHolder struct {
		MarkdownString string
	}
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	err := enc.Encode(markdownHolder{MarkdownString: l.markdownOnInspecting(th.Token)})
	//notify.Debugf("inspected:%s", th.Token)
	if err != nil {
		log.Printf("[melrose.error] %#v\n", err)
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
