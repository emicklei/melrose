package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/ui/img"
	"github.com/fogleman/gg"
)

// pianorollImageHandler returns PNG content given a variable
func (l *LanguageServer) pianorollImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		notify.Console.Warnf("HTTP method not allowed:%s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	varname := r.URL.Query().Get("var")
	if varname == "" {
		fmt.Fprintf(w, "var parameter is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	obj, ok := l.context.Variables().Get(varname)
	if !ok {
		fmt.Fprintf(w, "no object found for var:%s", varname)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	seq, ok := obj.(core.Sequenceable)
	if !ok {
		fmt.Fprintf(w, "no object is not sequenceable:%s", varname)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tim := core.NewTimeline()
	d := midi.NewOutputDevice(0, nil, 0, tim)
	d.Play(core.NoCondition, seq, l.context.Control().BPM(), time.Now())

	gc := gg.NewContext(500, 50)

	evts := tim.NoteEvents()
	nv := img.NotesView{Events: evts}
	nv.DrawOn(gc)

	w.Header().Set("content-type", "image/png")
	gc.EncodePNG(w)
}
