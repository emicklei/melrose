package control

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Knob struct {
	deviceID int
	channel  int
	number   int
	// set when used in assignment
	variableName string
	// changes
	currentValue int
}

func NewKnob(deviceID, channel, number int) *Knob {
	return &Knob{deviceID: deviceID, channel: channel, number: number}
}

// Inspect is part of Inspectable
func (k *Knob) Inspect(i core.Inspection) {
	i.Properties["device"] = k.deviceID
	i.Properties["channel"] = k.channel
	i.Properties["number"] = k.number
	i.Properties["currentValue"] = k.currentValue
}

// Storex is part of core.Storable
func (k *Knob) Storex() string {
	return fmt.Sprintf("knob(%d,%d)", k.deviceID, k.number)
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
	// if number != number {
	// 	return
	// }
	notify.Infof("knob %s (%d,%d,%d) = %d", k.variableName, k.deviceID, k.channel, k.number, value)
	k.currentValue = value
}

func (k *Knob) Value() interface{} {
	return k.currentValue
}

// VariableName is part of NameAware
func (k *Knob) VariableName(yours string) {
	k.variableName = yours
}
