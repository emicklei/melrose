package midi

import (
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type KeyTrigger struct {
	mutex   *sync.RWMutex
	playing bool
	count   int
	ctx     core.Context
	channel int
	note    core.Note
	fun     core.HasValue
}

func NewKeyTrigger(ctx core.Context, channel int, onNote core.Note, startStop core.HasValue) *KeyTrigger {
	return &KeyTrigger{
		mutex:   new(sync.RWMutex),
		ctx:     ctx,
		channel: channel,
		note:    onNote,
		fun:     startStop}
}

// NoteOn is part of core.NoteListener
func (t *KeyTrigger) NoteOn(channel int, n core.Note) {
	if notify.IsDebug() {
		notify.Debugf("keytrigger.NoteOn ch=%d note=%v", channel, n)
	}
	if channel != t.channel {
		return
	}
	if n.Name != t.note.Name {
		return
	}
	if n.Octave != t.note.Octave {
		return
	}
	val := t.fun.Value()
	if val == nil {
		return
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	// both playable and evaluatable are allowed
	if play, ok := val.(core.Playable); ok {
		if t.playing {
			notify.Infof("%s -> stop(%s)", t.note.String(), core.Storex(t.fun))
			t.playing = false
		} else {
			t.playing = true
			t.count++
			activeCount := t.count
			notify.Infof("%s -> play(%s)", t.note.String(), core.Storex(t.fun))
			_ = play.Play(t.ctx, func() bool { return t.playing && t.count == activeCount }, time.Now())
		}
	}
	// not playable, maybe evaluatable
	if eval, ok := val.(core.Evaluatable); ok {
		eval.Evaluate(t.ctx)
	}
}

// NoteOff is part of core.NoteListener
func (t *KeyTrigger) NoteOff(channel int, n core.Note) {
	if notify.IsDebug() {
		notify.Debugf("keytrigger.NoteOff ch=%d note=%v", channel, n)
	}
	// key trigger is not interested in this
}

func (t *KeyTrigger) ControlChange(channel, number, value int) {
	if notify.IsDebug() {
		notify.Debugf("keytrigger.ControlChange %d %d %d", channel, number, value)
	}
	// key trigger is not interested in this
}
