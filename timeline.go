package melrose

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// TODO use sync.Pool for noteEvents
type Timeline struct {
	head       *scheduledTimelineEvent // earliest
	tail       *scheduledTimelineEvent // latest
	protection sync.RWMutex
	resume     chan bool
	verbose    bool
}

func NewTimeline() *Timeline {
	return &Timeline{
		protection: sync.RWMutex{},
		resume:     make(chan bool),
		verbose:    false,
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
	if t.verbose {
		log.Println(event, when.Sub(now))
	}
	diff := when.Sub(now)
	// if between -wait..wait then handle now
	if -wait <= diff && diff <= wait {
		event.Handle(t, now)
		return nil
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
	if event.when.After(s.tail.when) {
		// event is after tail, new tail
		s.tail.next = event
		s.tail = event
		return
	}
	if s.head.when.After(event.when) {
		// event is before head, new head
		event.next = s.head
		s.head = event
		return
	}
	if s.head.next == nil {
		// event on the same time as head, put it after! head
		s.head.next = event
		s.tail = event
		return
	}
	if s.tail.when == event.when {
		// event on the same time as head, put it after! tail
		s.tail.next = event
		s.tail = event
		return
	}
	// somewhere between head and tail
	previous := s.head
	here := s.head.next
	for event.when.After(here.when) {
		previous = here
		here = here.next
	}
	// here is after event, it must be scheduled before it
	previous.next = event
	event.next = here
}
