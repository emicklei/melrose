package melrose

import (
	"fmt"
	"sync"
	"time"
)

// TODO use sync.Pool for noteEvents
type Timeline struct {
	head       *scheduledTimelineEvent // earliest
	tail       *scheduledTimelineEvent // latest
	protection sync.RWMutex
	resume     chan bool
}

func NewTimeline() *Timeline {
	return &Timeline{
		protection: sync.RWMutex{},
		resume:     make(chan bool),
	}
}

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

func (t *Timeline) Run() {
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
func (s *Timeline) Reset() {
	s.protection.Lock()
	defer s.protection.Unlock()
	s.head = nil
	s.tail = nil
}

func (t *Timeline) Schedule(event TimelineEvent, when time.Time) error {
	now := time.Now()
	if when.Before(now) {
		return fmt.Errorf("cannot schedule in the past:%v", now.Sub(when))
	}
	t.schedule(&scheduledTimelineEvent{
		when:  when,
		event: event,
	})
	return nil
}

// pre: event.when >= now
func (s *Timeline) schedule(event *scheduledTimelineEvent) {
	s.protection.Lock()
	if s.head == nil {
		s.head = event
		s.tail = event
		// before resume otherwise handle loop will deadlock
		s.protection.Unlock()
		s.resume <- true
		return
	}
	defer s.protection.Unlock()
	if s.head.when.After(event.when) {
		// call is before head, new head
		event.next = s.head
		s.head = event
		return
	}
	if event.when.After(s.tail.when) {
		// call is after tail, new tail
		s.tail.next = event
		s.tail = event
		return
	}
	if s.head.next == nil {
		// call on the same time as head, new head
		event.next = s.head
		s.head = event
		return
	}
	// somewhere between head and tail
	previous := s.head
	here := s.head.next
	for event.when.After(here.when) {
		previous = here
		here = here.next
	}
	// here is after call, it must be scheduled before it
	previous.next = event
	event.next = here
}
