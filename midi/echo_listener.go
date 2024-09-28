package midi

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

var DefaultEchoListener = EchoListener{}

type EchoListener struct {
}

func (e EchoListener) NoteOn(channel int, n core.Note) {}
func (e EchoListener) NoteOff(channel int, n core.Note) {
	fmt.Fprintf(notify.Console.StandardOut, "%s xx", n.String())
}
func (e EchoListener) ControlChange(channel, number, value int) {}
