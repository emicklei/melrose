package midi

import (
	"errors"
	"log"
	"math"
	"sync"
	"time"

	"github.com/emicklei/melrose"
	"github.com/rakyll/portmidi"
)

// TODO use sync.Pool for noteEvents
type Timeline struct {
	head       *noteEvent // earliest
	tail       *noteEvent // latest
	protection sync.RWMutex
	resume     chan bool
	out        *portmidi.Stream
}

func NewTimeline(out *portmidi.Stream) *Timeline {
	return &Timeline{
		protection: sync.RWMutex{},
		resume:     make(chan bool),
		out:        out,
	}
}

type TimelineEvent interface {
	Handle(tim *Timeline)
}

type scheduledTimelineEvent struct {
	event TimelineEvent
	when  time.Time
	next  *scheduledTimelineEvent
}

type noteEvent struct {
	when     time.Time
	which    []int64
	onoff    int
	channel  int
	velocity int64
	next     *noteEvent
}

func (e noteEvent) perform(out *portmidi.Stream) {
	//log.Println(s.which, s.onoff, s.channel)
	for _, each := range e.which {
		out.WriteShort(int64(e.onoff|(e.channel-1)), each, e.velocity)
	}
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
			//log.Println("wait for resume")
			<-t.resume
			continue
		}
		now := time.Now()
		for now.After(here.when) {
			here.perform(t.out)

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
				//log.Println("standard wait", wait)
				time.Sleep(wait) // 1/16 note
			} else {
				//log.Println("wait for next", untilNext)
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

// Schedule adds a Send to be performed in the future.
func (t *Timeline) Schedule(delayMs int64, which []int64, onoff int) error {
	absoluteTime := time.Now().Add(time.Duration(delayMs) * time.Millisecond)
	call := &noteEvent{
		when:     absoluteTime,
		which:    which,
		onoff:    onoff,
		channel:  1,
		velocity: 70,
	}
	return t.schedule(call, delayMs)
}

func (t *Timeline) add(when time.Time, which []int64, onoff int) error {
	call := &noteEvent{
		when:     when,
		which:    which,
		onoff:    onoff,
		channel:  1,
		velocity: 70,
	}
	return t.schedule(call, when.Sub(time.Now()).Milliseconds())
}

func (s *Timeline) schedule(call *noteEvent, delayMs int64) error {
	if delayMs < 0 {
		return errors.New("cannot schedule a call in the past")
	}
	s.protection.Lock()
	if s.head == nil {
		s.head = call
		s.tail = call
		// before resume otherwise send loop will deadlock
		s.protection.Unlock()
		s.resume <- true
		return nil
	}
	defer s.protection.Unlock()
	if s.head.when.After(call.when) {
		// call is before head, new head
		call.next = s.head
		s.head = call
		return nil
	}
	if call.when.After(s.tail.when) {
		// call is after tail, new tail
		s.tail.next = call
		s.tail = call
		return nil
	}
	if s.head.next == nil {
		// call on the same time as head, new head
		call.next = s.head
		s.head = call
		return nil
	}
	// somewhere between head and tail
	previous := s.head
	here := s.head.next
	for call.when.After(here.when) {
		previous = here
		here = here.next
	}
	// here is after call, it must be scheduled before it
	previous.next = call
	call.next = here

	return nil
}

// TEMP

// schedule all the notes on the timeline
func (t *Timeline) Play(s melrose.Sequenceable) {
	wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/120.))) * time.Millisecond
	actualSequence := s.S()
	moment := time.Now()
	for _, eachGroup := range actualSequence.Notes {
		var actualDuration time.Duration
		for _, eachNote := range eachGroup {
			// all have the same duration so combine the event
			actualDuration = time.Duration(float32(wholeNoteDuration) * eachNote.DurationFactor())
			if eachNote.IsRest() {
				continue
			}
			nr := int64(eachNote.MIDI())
			if err := t.add(moment, []int64{nr}, noteOn); err != nil {
				log.Println(err)
				return
			}
			if err := t.add(moment.Add(actualDuration), []int64{nr}, noteOff); err != nil {
				log.Println(err)
				return
			}
		}
		moment = moment.Add(actualDuration)
	}
}
