package melrose

import (
	"fmt"
	"sync"
	"time"
)

type Loop struct {
	Target    Sequenceable
	isRunning bool
	mutex     sync.RWMutex
}

func NewLoop(target Sequenceable) *Loop {
	return &Loop{
		Target: target,
	}
}

func (l *Loop) Storex() string {
	if s, ok := l.Target.(Storable); ok {
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
}

func (l *Loop) reschedule(d AudioDevice, when time.Time) {
	endOfLastNote := d.Play(l.Target, Context().LoopControl.BPM(), when)
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
	l.reschedule(Context().AudioDevice, when)
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
	l.Target = newTarget
}

// Play is part of Playable
func (l *Loop) Play(d AudioDevice) error {
	Context().LoopControl.Begin(l)
	return nil
}
