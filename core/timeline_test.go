package core

import (
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

type testEvent struct {
	id int
}

func (e testEvent) NoteChangesDo(block func(NoteChange)) {}
func (e testEvent) Handle(t *Timeline, w time.Time)      {}

func TestScheduleAdd(t *testing.T) {
	tim := NewTimeline()
	now := time.Now()

	e1 := testEvent{id: 1}
	e2 := testEvent{id: 2}
	e3 := testEvent{id: 3}
	e4 := testEvent{id: 4}
	// e1 -> e2 -> e4 -> e3
	tim.Schedule(e1, now.Add(1*time.Second))
	if got, want := tim.head.event, e1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tim.Schedule(e2, now.Add(1*time.Second))
	if got, want := tim.head.next.event, e2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tim.Schedule(e3, now.Add(5*time.Second))
	if got, want := tim.tail.event, e3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	tim.Schedule(e4, now.Add(3*time.Second))
	if got, want := tim.head.next.next.event, e4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.tail.event, e3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.head.next.next.next.event, e3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}

	zs := tim.ZeroStarting()
	if got, want := zs.head.when.Second(), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := zs.head.next.when.Second(), 0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := zs.head.next.next.when.Second(), 2; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := zs.head.next.next.next.when.Second(), 4; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
