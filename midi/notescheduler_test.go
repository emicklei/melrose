package midi

import (
	"log"
	"testing"
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/m"
)

// go test -v -run "^(TestNotescheduler)$"
func TestNotescheduler(t *testing.T) {
	t.Skip()
	dev, _ := Open()
	defer dev.Close()
	melrose.SetCurrentDevice(dev)
	dev.echo = true
	dev.bpm = 120

	tim := NewTimeline(dev.stream)
	go tim.Run()
	defer tim.Reset()

	tim.Play(m.Sequence("C E G B D F A C5 [C E G]"))

	time.Sleep(8 * time.Second)
}

func TestScheduleAddInThePast(t *testing.T) {
	tim := NewTimeline(nil)
	now := time.Now()
	call := &noteEvent{when: now}
	if err := tim.schedule(call, -1); err == nil {
		t.Fatal("error expected")
	}
}

func TestScheduleAdd(t *testing.T) {
	tim := NewTimeline(nil)
	go func() {
		for {
			log.Println(<-tim.resume)
		}
	}()
	now := time.Now()
	e1 := &noteEvent{when: now.Add(1 * time.Second)}
	e2 := &noteEvent{when: now.Add(1 * time.Second)}
	e3 := &noteEvent{when: now.Add(5 * time.Second)}
	e4 := &noteEvent{when: now.Add(3 * time.Second)}
	tim.schedule(e1, 1000)
	tim.schedule(e2, 1000)
	tim.schedule(e3, 5000)
	tim.schedule(e4, 3000)

	if got, want := tim.head, e2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.tail, e3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.head.next, e1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.head.next.next, e4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tim.head.next.next.next, e3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
