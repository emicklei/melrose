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
		w.WriteHeader(http.StatusMethodNotAllowed)
		notify.Console.Warnf("HTTP method not allowed:%s", r.Method)
		return
	}
	varname := r.URL.Query().Get("var")
	if varname == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "var parameter is empty")
		return
	}
	obj, ok := l.context.Variables().Get(varname)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "no object found for var:%s", varname)
		return
	}
	seq, ok := obj.(core.Sequenceable)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "no object is not sequenceable:%s", varname)
		return
	}

	tim := core.NewTimeline()
	d := midi.NewOutputDevice(0, nil, 0, tim)
	d.Play(core.NoCondition, seq, l.context.Control().BPM(), time.Now())

	gc := gg.NewContext(1000, 200)

	evts := tim.NoteEvents()
	nv := img.NotesView{Events: evts, BPM: l.context.Control().BPM()}
	nv.DrawOn(gc)

	w.Header().Set("content-type", "image/png")
	gc.EncodePNG(w)
}
