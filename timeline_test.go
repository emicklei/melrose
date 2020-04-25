package melrose

import (
	"log"
	"testing"
	"time"
)

func TestScheduleAddInThePast(t *testing.T) {
	tim := NewTimeline()
	now := time.Now()
	call := new(testEvent)
	if err := tim.Schedule(call, now.Add(-1*time.Second)); err == nil {
		t.Fatal("error expected")
	}
}

type testEvent struct{}

func (e testEvent) Handle(t *Timeline, w time.Time) {}

func TestScheduleAdd(t *testing.T) {
	tim := NewTimeline()
	go func() {
		for {
			log.Println(<-tim.resume)
		}
	}()
	now := time.Now()

	e1 := new(testEvent)
	e2 := new(testEvent)
	e3 := new(testEvent)
	e4 := new(testEvent)
	tim.Schedule(e1, now.Add(1*time.Second))
	tim.Schedule(e2, now.Add(1*time.Second))
	tim.Schedule(e3, now.Add(5*time.Second))
	tim.Schedule(e4, now.Add(3*time.Second))

	if got, want := tim.head.event, e2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.tail.event, e3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.head.next.event, e1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.head.next.next.event, e4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.head.next.next.next.event, e3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
