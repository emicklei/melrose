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
	ctx     core.Context // needed?
	note    core.Note    // or valuabLe?
	fun     core.Valueable
}

func NewKeyTrigger(ctx core.Context, onNote core.Note, startStop core.Valueable) *KeyTrigger {
	return &KeyTrigger{
		mutex: new(sync.RWMutex),
		ctx:   ctx,
		note:  onNote,
		fun:   startStop}
}

// NoteOn is part of core.NoteListener
func (t *KeyTrigger) NoteOn(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("keytrigger.NoteOn %v", n)
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
			play.Stop(t.ctx)
			t.playing = false
		} else {
			_ = play.Play(t.ctx, time.Now())
			t.playing = true
		}
		return
	}
	if eval, ok := val.(core.Evaluatable); ok {
		eval.Evaluate(t.ctx)
	}
}

// NoteOff is part of core.NoteListener
func (t *KeyTrigger) NoteOff(n core.Note) {
	if core.IsDebug() {
		notify.Debugf("keytrigger.NoteOff %v", n)
	}
	// key trigger is not interested in this
}
