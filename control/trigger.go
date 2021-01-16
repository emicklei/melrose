package control

import (
	"errors"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Trigger struct {
	ctx      core.Context
	deviceID int
	note     core.Note // or valuabLe?
	fun      core.Valueable
}

func NewTrigger(ctx core.Context, deviceID int, onNote core.Note, startStop core.Valueable) *Trigger {
	return &Trigger{
		ctx:      ctx,
		deviceID: deviceID,
		note:     onNote,
		fun:      startStop}
}

// Play is part of core.Playable
func (t *Trigger) Play(ctx core.Context, at time.Time) error {
	// ignore time for now
	if !ctx.Device().HasInputCapability() {
		return errors.New("Input is not available for this device")
	}
	ctx.Device().Listen(t.deviceID, t, true)
	return nil
}

// Stop is part of core.Playable
func (t *Trigger) Stop(ctx core.Context) error {
	return nil
}

// NoteOn is part of core.NoteListener
func (t *Trigger) NoteOn(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("trigger.NoteOn %v", n)
	}
}

// NoteOff is part of core.NoteListener
func (t *Trigger) NoteOff(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("trigger.NoteOff %v", n)
	}
}
