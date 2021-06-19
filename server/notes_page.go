package server

import (
	"fmt"
	"net/http"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

// notesPageHandler returns HTML content given a variable
func (l *LanguageServer) notesPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		notify.Console.Warnf("HTTP method not allowed:%s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	varname := r.URL.Query().Get("var")

	// temp
	if varname == "" {
		l.context.Variables().Put("test", core.MustParseSequence("c d e f g"))
		varname = "test"
	}

	// object can refer to one or more devices
	// object can refer to one or more channels per device
	// create notes view for each device,channel pair
	fmt.Fprintf(w, `
	<html>
		<body>
			<h1>%s</h1>
			<img src="/v1/pianoroll?var=%s"></img>
		</body>		
	</html>

	`, varname, varname)
}
