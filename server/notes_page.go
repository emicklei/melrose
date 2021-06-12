package server

import (
	"net/http"

	"github.com/emicklei/melrose/notify"
)

// notesPageHandler returns HTML content given a variable
func (l *LanguageServer) notesPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		notify.Console.Warnf("HTTP method not allowed:%s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	notify.Println(r.URL.Query().Get("var"))

	// object can refer to one or more devices
	// object can refer to one or more channels per device
	// create notes view for each device,channel pair
}
