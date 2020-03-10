package main

import (
	. "github.com/emicklei/melrose"
	"github.com/emicklei/melrose/pilot"
)

// go run scale_pilot.go
func main() {
	p := pilot.Open()
	defer p.Close()

	// https://github.com/hundredrabbits/Pilot/blob/master/desktop/sources/scripts/mixer.js
	// p.Send("1OSCsisq")
	// p.Send("reset")
	// p.Send("rosc")
	cm := C().Scale(Major)

	p.Play(cm.S())
}
