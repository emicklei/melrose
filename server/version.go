package server

import (
	"encoding/json"
	"net/http"

	"github.com/emicklei/melrose/core"
)

type versionInfo struct {
	APIVersion string
	BuildTag   string
}

func (l *LanguageServer) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("access-control-allow-origin", "*")
	v := versionInfo{
		APIVersion: "v1",
		BuildTag:   core.BuildTag,
	}
	json.NewEncoder(w).Encode(v)
}
