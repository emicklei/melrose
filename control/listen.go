package control

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Listen struct {
	deviceID      int
	variableStore core.VariableStorage
	variableName  string
	isRunning     bool
	callback      core.Valueable
}

func NewListen(deviceID int, store core.VariableStorage, variableName string, target core.Valueable) *Listen {
	return &Listen{
		deviceID:      deviceID,
		variableStore: store,
		variableName:  variableName,
		callback:      target,
	}
}

// Play is part of core.Playable
func (l *Listen) Play(ctx core.Context) error {
	if l.isRunning {
		return nil
	}
	l.isRunning = true
	ctx.Device().Listen(l.deviceID, l, l.isRunning)
	return nil
}

func (l *Listen) Stop(ctx core.Context) {
	if l.isRunning {
		return
	}
	l.isRunning = false
	ctx.Device().Listen(l.deviceID, l, l.isRunning)
}

// NoteOn is part of core.NoteListener
func (l *Listen) NoteOn(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("control.listen ON %v", n)
	}
	l.variableStore.Put(l.variableName, n)
	if e, ok := l.callback.Value().(core.Evaluatable); ok {
		e.Evaluate()
	}
}

// NoteOff is part of core.NoteListener
func (l *Listen) NoteOff(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("control.listen OFF %v", n)
	}
}

// Storex is part of core.Storable
func (l *Listen) Storex() string {
	return fmt.Sprintf("listen(%d,%s,%s)", l.deviceID, l.variableName, core.Storex(l.callback))
}
