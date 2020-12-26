package core

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

type Loop struct {
	ctx        Context
	target     []Sequenceable
	isRunning  bool
	mutex      sync.RWMutex
	condition  Condition
	nextPlayAt time.Time
}

func NewLoop(ctx Context, target []Sequenceable) *Loop {
	return &Loop{
		ctx:       ctx,
		target:    target,
		condition: TrueCondition,
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

func (l *Loop) Evaluate(condition Condition) error {
	// create a start a clone
	clone := NewLoop(l.ctx, l.target)
	clone.condition = condition
	if IsDebug() {
		notify.Debugf("loop.eval")
	}
	clone.Start(l.ctx.Device())
	return nil
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

// in mutex
func (l *Loop) reschedule(d AudioDevice, when time.Time) {
	// check condition first
	if l.condition != nil && !l.condition() {
		l.isRunning = false
		return
	}
	moment := when
	for _, each := range l.target {
		// after each other
		moment = d.Play(l.condition, each, l.ctx.Control().BPM(), moment)
	}
	if IsDebug() {
		notify.Debugf("core.loop: next=%s", moment.Format("15:04:05.00"))
	}
	// schedule the loop itself so it can play again when Handle is called
	l.nextPlayAt = moment
	d.Schedule(l, moment)
}

func (l *Loop) NextPlayAt() time.Time {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if !l.isRunning {
		return time.Time{}
	}
	return l.nextPlayAt
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

type PlaySynchronizer struct {
	players []Valueable
}

func NewPlaySynchronizer(vals []Valueable) PlaySynchronizer {
	return PlaySynchronizer{players: vals}
}

// Play is part of Playable
func (s PlaySynchronizer) Play(ctx Context) error {
	notify.Debugf("playsync.play")
	// first for loops
	if len(s.players) == 0 {
		return nil
	}
	first, ok := s.players[0].Value().(*Loop)
	if !ok {
		return nil
	}
	if !first.isRunning {
		return nil
	}
	// when := first.NextPlayAt()
	// for i:=1;i<len(s.players);i++ {
	// 	p, ok  := s.players[i].Value().(*Loop)
	// 	if ok {
	// 		p.Play(when)
	// 	}
	// }
	return nil
}

func (s PlaySynchronizer) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "onnext(")
	for i, each := range s.players {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, "%s", Storex(each))
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}
