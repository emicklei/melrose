package core

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

type Loop struct {
	ctx       Context
	target    []Sequenceable
	isRunning bool
	mutex     sync.RWMutex
}

func NewLoop(ctx Context, target []Sequenceable) *Loop {
	return &Loop{
		ctx:    ctx,
		target: target,
	}
}

func (l *Loop) Target() []Sequenceable { return l.target }

func (l *Loop) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "loop(")
	AppendStorexList(&b, true, l.target)
	fmt.Fprintf(&b, ")")
	return b.String()
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
	moment := when
	for _, each := range l.target {
		// after each other
		moment = d.Play(NoCondition, each, l.ctx.Control().BPM(), moment)
	}
	if IsDebug() {
		notify.Debugf("core.loop: next=%s", moment.Format("15:04:05.00"))
	}
	// schedule the loop itself so it can play again when Handle is called
	d.Schedule(l, moment)
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

func (l *Loop) SetTarget(newTarget []Sequenceable) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.target = newTarget
}

// Play is part of Playable
func (l *Loop) Play(ctx Context) error {
	ctx.Control().StartLoop(l)
	return nil
}
