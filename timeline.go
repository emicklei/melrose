package melrose

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Timeline is a chain of events that are place in the future (playing) or past (recording).
// TODO use sync.Pool for noteEvents
type Timeline struct {
	head       *scheduledTimelineEvent // earliest
	tail       *scheduledTimelineEvent // latest
	protection sync.RWMutex
	isPlaying  bool
	resume     chan bool
	verbose    bool
}

// NewTimeline creates a new Timeline.
func NewTimeline() *Timeline {
	return &Timeline{
		protection: sync.RWMutex{},
		verbose:    false,
	}
}

// TimelineEvent describe an event that can be scheduled on a Timeline.
type TimelineEvent interface {
	Handle(tim *Timeline, when time.Time)
}

type scheduledTimelineEvent struct {
	event TimelineEvent
	when  time.Time
	next  *scheduledTimelineEvent
}

var (
	wait = 50 * time.Millisecond // 1/16 note @ bpm 300
)

// Len returns the current number of scheduled events.
func (t *Timeline) Len() int64 {
	t.protection.RLock()
	defer t.protection.RUnlock()
	var count int64
	here := t.head
	for here != nil {
		here = here.next
		count++
	}
	return count
}

// IsEmpty is true if no events are scheduled.
func (t *Timeline) IsEmpty() bool {
	t.protection.RLock()
	defer t.protection.RUnlock()
	return t.head == nil
}

// Play runs a loop to handle all the events in time. This is blocking.
func (t *Timeline) Play() {
	t.resume = make(chan bool)
	t.isPlaying = true
	for {
		t.protection.RLock()
		here := t.head
		t.protection.RUnlock()
		if here == nil {
			<-t.resume
			continue
		}
		now := time.Now()
		for now.After(here.when) {
			here.event.Handle(t, now)

			t.protection.Lock()
			t.head = t.head.next
			here = t.head
			t.protection.Unlock()

			if here == nil {
				break
			}
		}
		if here != nil {
			untilNext := here.when.Sub(now)
			if wait < untilNext {
				time.Sleep(wait) // 1/16 note
			} else {
				time.Sleep(untilNext) // < 1/16 note
			}
		}
	}
}

// Reset forgets about all scheduled calls.
func (t *Timeline) Reset() {
	t.protection.Lock()
	defer t.protection.Unlock()
	t.head = nil
	t.tail = nil
}

// Schedule adds an event for a given time
func (t *Timeline) Schedule(event TimelineEvent, when time.Time) error {
	now := time.Now()
	if t.verbose {
		log.Println(event, when.Sub(now))
	}
	diff := when.Sub(now)
	if t.isPlaying {
		// if between -wait..wait then handle now
		if -wait <= diff && diff <= wait {
			event.Handle(t, now)
			return nil
		}
	}
	if diff < -wait {
		return fmt.Errorf("cannot schedule in the past:%v", now.Sub(when))
	}
	t.schedule(&scheduledTimelineEvent{
		when:  when,
		event: event,
	})
	return nil
}

// schedule adds an event on the chain.
// pre: event.when >= now
func (t *Timeline) schedule(event *scheduledTimelineEvent) {
	t.protection.Lock()
	if t.head == nil {
		t.head = event
		t.tail = event
		// before resume otherwise run loop will deadlock
		t.protection.Unlock()
		if t.isPlaying {
			t.resume <- true
		}
		return
	}
	defer t.protection.Unlock()
	if event.when.After(t.tail.when) {
		// event is after tail, new tail
		t.tail.next = event
		t.tail = event
		return
	}
	if t.head.when.After(event.when) {
		// event is before head, new head
		event.next = t.head
		t.head = event
		return
	}
	if t.head.next == nil {
		// event on the same time as head, put it after! head
		t.head.next = event
		t.tail = event
		return
	}
	if t.tail.when == event.when {
		// event on the same time as head, put it after! tail
		t.tail.next = event
		t.tail = event
		return
	}
	// somewhere between head and tail
	previous := t.head
	here := t.head.next
	for event.when.After(here.when) {
		previous = here
		here = here.next
	}
	// here is after event, it must be scheduled before it
	previous.next = event
	event.next = here
}

// EventsDo visits all scheduled events and call the block for each.
func (t *Timeline) EventsDo(block func(event TimelineEvent, when time.Time)) {
	t.protection.Lock()
	defer t.protection.Unlock()
	here := t.head
	for here != nil {
		block(here.event, here.when)
		here = here.next
	}
}
