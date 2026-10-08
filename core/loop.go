package core

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

// todo: protect with mutex
var runningLoop *Loop

type Loop struct {
	ctx        Context
	target     []Sequenceable
	isRunning  bool
	mutex      sync.RWMutex
	condition  Condition
	startedAt  time.Time
	nextPlayAt time.Time
	clock      *PlaybackClock
	cycleTicks int64
	cycleFixed time.Duration
}

func NewLoop(ctx Context, target []Sequenceable) *Loop {
	return &Loop{
		ctx:       ctx,
		target:    target,
		condition: TrueCondition,
	}
}

func (l *Loop) Target() []Sequenceable { return l.target }

func (l *Loop) SetTarget(newTarget []Sequenceable) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.target = newTarget
}

func (l *Loop) IsRunning() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.isRunning
}

func (l *Loop) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "loop(")
	AppendStorexList(&b, true, l.target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (l *Loop) Evaluate(ctx Context) error {
	// create and start a clone
	clone := NewLoop(l.ctx, l.target)
	cond := NoCondition
	if with, ok := ctx.(Conditional); ok {
		cond = with.Condition()
	}
	clone.condition = cond
	if notify.IsDebug() {
		notify.Debugf("loop.eval")
	}
	clone.Play(l.ctx, cond, time.Now())
	return nil
}

// Inspect is part of Inspectable
func (l *Loop) Inspect(i Inspection) {
	i.Properties["running"] = l.isRunning
	if runningLoop == l {
		i.Properties["leader"] = true
	}
}

// in mutex
func (l *Loop) reschedule(d AudioDevice) {
	if !l.isRunning {
		return
	}
	if l.condition != nil && !l.condition() {
		l.isRunning = false
		return
	}
	bpm := l.ctx.Control().BPM()
	l.clock.SetBPM(bpm)
	moment := l.clock.Time()
	begin := moment
	startClock := *l.clock
	clocked, usesClock := d.(ClockedAudioDevice)
	for _, each := range l.target {
		// after each other
		if usesClock {
			moment = clocked.PlayWithClock(l.condition, each, l.clock)
		} else {
			moment = d.Play(l.condition, each, bpm, moment)
			*l.clock = NewPlaybackClock(moment, bpm)
		}
	}
	if notify.IsDebug() {
		notify.Debugf("core.loop: next=%s", moment.Format("15:04:05.00"))
	}
	if !moment.After(begin) {
		l.isRunning = false
		if runningLoop == l {
			runningLoop = nil
		}
		return
	}
	if usesClock {
		l.cycleTicks = l.clock.ticks - startClock.ticks
		l.cycleFixed = l.clock.fixed - startClock.fixed
	} else {
		l.cycleTicks = 0
		l.cycleFixed = moment.Sub(begin)
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
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.isRunning {
		return
	}
	l.clock.SetBPM(l.ctx.Control().BPM())
	period := PlaybackClock{bpm: l.clock.bpm, ticks: l.cycleTicks, fixed: l.cycleFixed}.Duration()
	next := *l.clock
	next.ticks += l.cycleTicks
	next.fixed += l.cycleFixed
	if period > 0 && !next.Time().After(when) {
		missed := int64(when.Sub(l.clock.Time()) / period)
		if missed < 1 {
			missed = 1
		}
		l.clock.ticks += missed * l.cycleTicks
		l.clock.fixed += time.Duration(missed) * l.cycleFixed
		for l.clock.Time().After(when) {
			l.clock.ticks -= l.cycleTicks
			l.clock.fixed -= l.cycleFixed
		}
		for {
			end := *l.clock
			end.ticks += l.cycleTicks
			end.fixed += l.cycleFixed
			if end.Time().After(when) {
				break
			}
			*l.clock = end
		}
	}
	l.reschedule(l.ctx.Device())
}

func (l *Loop) NoteChangesDo(block func(NoteChange)) {}

// Play is part of Playable
func (l *Loop) Play(ctx Context, while Condition, at time.Time) time.Time {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	forever := time.Now().AddDate(100, 0, 0)
	if l.isRunning {
		return forever
	}
	when := at
	if runningLoop != nil {
		// only if loops do want to start at the same time
		// we delay this loop until the last loop finished
		if at.Sub(runningLoop.startedAt).Milliseconds() > 100 {
			when = runningLoop.nextPlayAt
		}
	} else {
		runningLoop = l
	}
	l.isRunning = true
	l.condition = while
	l.startedAt = when
	clock := NewPlaybackClock(when, ctx.Control().BPM())
	l.clock = &clock
	l.reschedule(l.ctx.Device())
	return forever
}

// Stop is part of Playable
func (l *Loop) Stop(ctx Context) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if !l.isRunning {
		return nil
	}
	l.isRunning = false

	if l == runningLoop {
		runningLoop = nil
	}

	return nil
}

// IsPlaying is part of Playable
func (l *Loop) IsPlaying() bool {
	return l.isRunning
}

func (l *Loop) S() Sequence {
	return l.ToSequence(1)
}

func (l *Loop) ToSequence(loopcount int) Sequence {
	all := Sequence{}
	for i := 0; i < loopcount; i++ {
		for _, each := range l.target {
			all = all.SequenceJoin(each.S())
		}
	}
	return all
}
