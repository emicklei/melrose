package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

type Loop struct {
	ctx       Context
	target    Sequenceable
	isRunning bool
	mutex     sync.RWMutex
}

func NewLoop(ctx Context, target Sequenceable) *Loop {
	return &Loop{
		ctx:    ctx,
		target: target,
	}
}

func (l *Loop) Target() Sequenceable { return l.target }

func (l *Loop) Storex() string {
	if s, ok := l.target.(Storable); ok {
		return fmt.Sprintf("loop(%s)", s.Storex())
	}
	return ""

}

func (l *Loop) IsRunning() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.isRunning
}

func (l *Loop) Start(d AudioDevice) *Loop {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.isRunning || d == nil {
		return l
	}
	l.isRunning = true
	l.reschedule(d, time.Now())
	return l
}

func (l *Loop) Inspect(i Inspection) {
	i.Properties["running"] = l.isRunning
	if st, ok := l.target.(Storable); ok {
		i.Properties["sequence"] = st.Storex()
	}
}

func (l *Loop) reschedule(d AudioDevice, when time.Time) {
	endOfLastNote := d.Play(l.target, l.ctx.Control().BPM(), when)
	if IsDebug() {
		notify.Debugf("next loop until [%s]", endOfLastNote.Format("15:04:05.00"))
	}
	// schedule the loop itself so it can play again when Handle is called
	d.Timeline().Schedule(l, endOfLastNote)
}

// Handle is part of TimelineEvent
func (l *Loop) Handle(tim *Timeline, when time.Time) {
	l.mutex.RLock()
	if !l.isRunning {
		l.mutex.RUnlock()
		return
	}
	l.mutex.RUnlock()
	l.reschedule(l.ctx.Device(), when)
}

func (l *Loop) Stop() *Loop {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if !l.isRunning {
		return l
	}
	l.isRunning = false
	return l
}

func (l *Loop) SetTarget(newTarget Sequenceable) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.target = newTarget
}

// Play is part of Playable
func (l *Loop) Play(ctx Context) error {
	ctx.Control().StartLoop(l)
	return nil
}
