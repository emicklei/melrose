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
	loopCount  int // if > 0 then stop after that many cycles; for tests
	cycles     int
}

func NewLoop(ctx Context, target []Sequenceable) *Loop {
	return &Loop{
		ctx:       ctx,
		target:    target,
		condition: TrueCondition,
	}
}

func (l *Loop) Target() []Sequenceable { return l.target }

// SetLoopCount makes the loop stop by itself after n cycles. Zero means forever.
func (l *Loop) SetLoopCount(n int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.loopCount = n
}

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
	clone.loopCount = l.loopCount
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

// reschedule queues one cycle from the planned playback boundary, not the current
// wall-clock time, so callback lateness does not accumulate as drift between loops.
// The caller must hold l.mutex and initialize l.clock before calling.
func (l *Loop) reschedule(d AudioDevice) {
	if !l.isRunning {
		return
	}
	if l.condition != nil && !l.condition() {
		l.isRunning = false
		return
	}
	// Apply tempo changes at the cycle boundary. SetBPM preserves that boundary
	// while making subsequent musical durations use the new tempo.
	bpm := l.ctx.Control().BPM()
	l.clock.SetBPM(bpm)
	moment := l.clock.Time()
	begin := moment
	startClock := *l.clock
	for _, each := range l.target {
		// A shared clock places targets consecutively and retains musical ticks
		// across cycles, avoiding repeated rounding of each sequence's duration.
		moment = d.PlayWithClock(l.condition, each, l.clock)
	}
	if notify.IsDebug() {
		notify.Debugf("core.loop: next=%s", moment.Format("15:04:05.00"))
	}
	// A cycle that advances no time would continually reschedule itself at the
	// same boundary, preventing the timeline from making progress.
	if !moment.After(begin) {
		l.isRunning = false
		if runningLoop == l {
			runningLoop = nil
		}
		return
	}
	// Remember this cycle's length so Handle can skip completed cycles after a
	// late wakeup without replaying or reevaluating every missed iteration.
	// Keep musical ticks separate from explicit durations: only ticks scale with BPM.
	l.cycleTicks = l.clock.ticks - startClock.ticks
	l.cycleFixed = l.clock.fixed - startClock.fixed
	l.cycles++
	if l.loopCount > 0 && l.cycles >= l.loopCount {
		l.isRunning = false
		if runningLoop == l {
			runningLoop = nil
		}
		return
	}
	// Queue the next cycle at this cycle's planned end, regardless of how long
	// scheduling took. Handle will resume from the clock already at that boundary.
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
	// Skip only when an entire additional cycle has elapsed. Smaller delays keep
	// the original boundary, so independently delayed loops retain their phase.
	if period > 0 && !next.Time().After(when) {
		missed := int64(when.Sub(l.clock.Time()) / period)
		if missed < 1 {
			missed = 1
		}
		l.clock.ticks += missed * l.cycleTicks
		l.clock.fixed += time.Duration(missed) * l.cycleFixed
		// Division by the rounded period estimates the skip count. Correct it
		// against absolute clock boundaries, whose cumulative rounding may differ,
		// to select the cycle containing when and schedule its end in the future.
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
	l.cycles = 0
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
