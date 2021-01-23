package control

import (
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Knob struct {
	vars     core.VariableStorage
	deviceID int
	channel  int
	number   int
	// changes
	currentValue int
}

func NewKnob(ctx core.Context, deviceID, channel, number int) *Knob {
	k := &Knob{vars: ctx.Variables(), deviceID: deviceID, channel: channel, number: number}
	ctx.Device().Listen(deviceID, k, true)
	return k
}

func (k *Knob) NoteOn(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("knob.NoteOn %v", n)
	}
}
func (k *Knob) NoteOff(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("knob.NoteOff %v", n)
	}
}
func (k *Knob) ControlChange(channel, number, value int) {
	if core.IsDebug() {
		notify.Debugf("knob.ControlChange ch=%d,nr=%d,val=%d", channel, number, value)
	}
	// TODO check channel
	if number != number {
		return
	}
	k.currentValue = value
}

func (k *Knob) Value() interface{} {
	return k.currentValue
}
