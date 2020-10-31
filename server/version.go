package server

import (
	"encoding/json"
	"net/http"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl"
)

type versionInfo struct {
	APIVersion    string
	SyntaxVersion string
	BuildTag      string
}

func (l *LanguageServer) versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	v := versionInfo{
		APIVersion:    "v1",
		SyntaxVersion: dsl.SyntaxVersion,
		BuildTag:      core.BuildTag,
	}
	json.NewEncoder(w).Encode(v)
	return
}
