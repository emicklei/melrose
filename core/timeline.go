package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

// Timeline is a chain of events that are placed in the future (playing).
type Timeline struct {
	head       *scheduledTimelineEvent // earliest
	tail       *scheduledTimelineEvent // latest
	protection sync.RWMutex
	isPlaying  bool
	resume     chan bool
}

// NewTimeline creates a new Timeline.
func NewTimeline() *Timeline {
	return &Timeline{
		protection: sync.RWMutex{},
	}
}

// TimelineEvent describes an event that can be scheduled on a Timeline.
type TimelineEvent interface {
	Handle(tim *Timeline, when time.Time)
	NoteChangesDo(block func(NoteChange))
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
	if notify.IsDebug() {
		notify.Debugf("core.timeline: flushing all scheduled MIDI events")
	}
	t.protection.Lock()
	defer t.protection.Unlock()
	t.head = nil
	t.tail = nil
}

// Schedule adds an event for a given time
func (t *Timeline) Schedule(event TimelineEvent, when time.Time) error {
	now := time.Now()
	diff := when.Sub(now)
	if diff < -wait {
		return fmt.Errorf("core.timeline: cannot schedule in the past:%v", now.Sub(when))
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
	if t.tail.when.Equal(event.when) {
		// event on the same time as tail, put it after! tail
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

// EventsDo visits all scheduled events and calls the block for each.
func (t *Timeline) EventsDo(block func(event TimelineEvent, when time.Time)) {
	t.protection.Lock()
	defer t.protection.Unlock()
	here := t.head
	for here != nil {
		block(here.event, here.when)
		here = here.next
	}
}

// ZeroStarting returns a new one in which all events are shifted back in time starting at time 0.
func (t *Timeline) ZeroStarting() *Timeline {
	if t.Len() == 0 {
		return t
	}
	result := NewTimeline()
	zero := time.Time{}
	t.EventsDo(func(event TimelineEvent, when time.Time) {
		d := when.Sub(t.head.when)
		result.schedule(&scheduledTimelineEvent{
			when:  zero.Add(d),
			event: event,
		})
	})
	return result
}

func (t *Timeline) NoteEvents() (list []NoteEvent) {
	activeNotes := map[int64]NoteEvent{}
	t.EventsDo(func(event TimelineEvent, when time.Time) {
		event.NoteChangesDo(func(change NoteChange) {
			if change.isOn {
				_, ok := activeNotes[change.note]
				if ok {
					// note was on ?
					// TODO warn?
				} else {
					// new
					activeNotes[change.note] = NoteEvent{Start: when, Number: change.Number(), Velocity: change.Velocity()}
				}
			} else {
				// note off
				hit, ok := activeNotes[change.note]
				if !ok {
					// note was never on ?
					// TODO warn?
				} else {
					list = append(list, hit.WithEnd(when))
					delete(activeNotes, change.note)
				}
			}
		})
	})
	return
}
