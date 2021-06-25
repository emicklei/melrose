package control

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi"
	"github.com/emicklei/melrose/notify"
)

const (
	noteOn  int64 = 0x90 // 10010000 , 144
	noteOff int64 = 0x80 // 10000000 , 128
)

type OnOff struct {
	isOn     bool
	deviceID int
	channel  int
	note     core.Note
}

func NewOnOff(deviceID int, channel int, note core.Note) OnOff {
	return OnOff{
		isOn:     false,
		deviceID: deviceID,
		channel:  channel,
		note:     note,
	}
}

// Play implements Playable
func (o OnOff) Play(ctx core.Context, at time.Time) error {
	if core.IsDebug() {
		notify.Debugf("control.OnOff.Play dev=%d ch=%d note=%v", o.deviceID, o.channel, o.note)
	}
	return o.send(ctx, noteOn)
}

func (o OnOff) send(ctx core.Context, status int64) error {
	nr := o.note.MIDI()
	velocity := o.note.Velocity
	mm := midi.NewMessage(
		ctx.Device(),
		core.On(o.deviceID),
		int(status),
		core.On(o.channel),
		core.On(nr),
		core.On(velocity))
	return mm.Evaluate(ctx)
}

// Stop implements Playable
func (o OnOff) Stop(ctx core.Context) error {
	if core.IsDebug() {
		notify.Debugf("control.OnOff.Stop dev=%d ch=%d note=%v", o.deviceID, o.channel, o.note)
	}
	return o.send(ctx, noteOff)
}

// Storex is part of core.Storable
func (o OnOff) Storex() string {
	return fmt.Sprintf("onoff(device(%d,channel(%d,%s)))", o.deviceID, o.channel, core.Storex(o.note))
}
