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

type timelinePlaybackEvent struct {
	testEvent
	handle func(time.Time)
}

func (e timelinePlaybackEvent) Handle(tim *Timeline, when time.Time) {
	e.handle(when)
}

func TestTimelinePlaybackRefreshesCallbackTime(t *testing.T) {
	timeline := NewTimeline()
	first := make(chan time.Time, 1)
	second := make(chan time.Time, 1)
	release := make(chan struct{})
	begin := time.Now().Add(10 * time.Millisecond)
	if err := timeline.Schedule(timelinePlaybackEvent{handle: func(when time.Time) {
		first <- when
		<-release
	}}, begin); err != nil {
		t.Fatal(err)
	}
	if err := timeline.Schedule(timelinePlaybackEvent{handle: func(when time.Time) {
		second <- when
	}}, begin); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		timeline.Play()
		close(done)
	}()
	t.Cleanup(func() {
		timeline.Stop()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Error("timeline did not stop")
		}
	})
	var firstTime time.Time
	select {
	case firstTime = <-first:
		close(release)
	case <-time.After(time.Second):
		close(release)
		t.Fatal("first callback did not run")
	}
	select {
	case secondTime := <-second:
		if !secondTime.After(firstTime) {
			t.Fatal("second callback received a stale execution timestamp")
		}
	case <-time.After(time.Second):
		t.Fatal("second callback did not run")
	}
}

type loopTimingDevice struct {
	AudioDeviceMock
}

type loopTimingSequence struct {
	Sequence
	calls int
}

func (s *loopTimingSequence) S() Sequence {
	s.calls++
	return s.Sequence
}

func (d *loopTimingDevice) Play(condition Condition, seq Sequenceable, bpm float64, beginAt time.Time) time.Time {
	return beginAt.Add(seq.S().DurationAt(bpm))
}

func (d *loopTimingDevice) PlayWithClock(condition Condition, seq Sequenceable, clock *PlaybackClock) time.Time {
	*clock = clock.AfterSequence(seq.S())
	return clock.Time()
}

func TestParallelLoopTiming(t *testing.T) {
	leader := runningLoop
	runningLoop = nil
	t.Cleanup(func() { runningLoop = leader })
	ctx := PlayContext{LoopControl: NewBeatmaster(nil, 123), AudioDevice: &loopTimingDevice{}}
	target := &loopTimingSequence{Sequence: S("1C")}
	first := NewLoop(ctx, []Sequenceable{target})
	second := NewLoop(ctx, []Sequenceable{S("C C C C")})
	begin := time.Now().Add(time.Hour)
	first.Play(ctx, NoCondition, begin)
	second.Play(ctx, NoCondition, begin)
	t.Cleanup(func() {
		first.Stop(ctx)
		second.Stop(ctx)
	})
	if !first.NextPlayAt().Equal(second.NextPlayAt()) {
		t.Fatal("equal-duration loops should initially share a boundary")
	}
	const iterations = 100
	for iteration := 0; iteration < iterations; iteration++ {
		first.Handle(nil, first.NextPlayAt().Add(time.Millisecond))
		second.Handle(nil, second.NextPlayAt().Add(3*time.Millisecond))
	}
	if got, want := second.NextPlayAt().Sub(first.NextPlayAt()), time.Duration(0); got != want {
		t.Errorf("callback lateness drift = %s, want %s", got, want)
	}
	if got, want := first.NextPlayAt().Sub(begin), fractionDuration(iterations+1, 123); got != want {
		t.Errorf("loop boundary = %s, want %s", got, want)
	}
	late := begin.Add(20 * time.Minute)
	first.Handle(nil, late)
	second.Handle(nil, late.Add(time.Millisecond))
	if !first.NextPlayAt().Equal(second.NextPlayAt()) || !first.NextPlayAt().After(late) {
		t.Fatal("overdue loops should skip completed iterations and resume on the same future boundary")
	}
	if got, want := target.calls, iterations+2; got != want {
		t.Errorf("sequence evaluated %d times, want %d played iterations", got, want)
	}
}

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

// BenchmarkSchedule benchmarks the schedule method with different scenarios
func BenchmarkSchedule(b *testing.B) {
	now := time.Now()

	b.Run("EmptyTimeline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tim := NewTimeline()
			event := &scheduledTimelineEvent{
				when:  now.Add(time.Duration(i) * time.Millisecond),
				event: testEvent{id: i},
			}
			tim.schedule(event)
		}
	})

	b.Run("AppendToTail", func(b *testing.B) {
		tim := NewTimeline()
		// Pre-populate with some events
		for i := 0; i < 100; i++ {
			tim.schedule(&scheduledTimelineEvent{
				when:  now.Add(time.Duration(i) * time.Millisecond),
				event: testEvent{id: i},
			})
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			event := &scheduledTimelineEvent{
				when:  now.Add(time.Duration(100+i) * time.Millisecond),
				event: testEvent{id: 100 + i},
			}
			tim.schedule(event)
		}
	})

	b.Run("PrependToHead", func(b *testing.B) {
		tim := NewTimeline()
		// Pre-populate with some events starting from a future time
		baseTime := now.Add(1 * time.Second)
		for i := 0; i < 100; i++ {
			tim.schedule(&scheduledTimelineEvent{
				when:  baseTime.Add(time.Duration(i) * time.Millisecond),
				event: testEvent{id: i},
			})
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			event := &scheduledTimelineEvent{
				when:  now.Add(time.Duration(i) * time.Millisecond),
				event: testEvent{id: -i},
			}
			tim.schedule(event)
		}
	})

	b.Run("InsertMiddle", func(b *testing.B) {
		tim := NewTimeline()
		// Pre-populate with events at regular intervals
		for i := 0; i < 100; i++ {
			tim.schedule(&scheduledTimelineEvent{
				when:  now.Add(time.Duration(i*10) * time.Millisecond),
				event: testEvent{id: i * 10},
			})
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Insert events in the middle of existing intervals
			event := &scheduledTimelineEvent{
				when:  now.Add(time.Duration(i%100*10+5) * time.Millisecond),
				event: testEvent{id: i%100*10 + 5},
			}
			tim.schedule(event)
		}
	})

	b.Run("LargeTimeline", func(b *testing.B) {
		tim := NewTimeline()
		// Pre-populate with many events
		for i := 0; i < 1000; i++ {
			tim.schedule(&scheduledTimelineEvent{
				when:  now.Add(time.Duration(i) * time.Millisecond),
				event: testEvent{id: i},
			})
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Insert events at random positions in the timeline
			event := &scheduledTimelineEvent{
				when:  now.Add(time.Duration(i%1000) * time.Millisecond),
				event: testEvent{id: 1000 + i},
			}
			tim.schedule(event)
		}
	})
}
