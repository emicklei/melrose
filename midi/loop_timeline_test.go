package midi

import (
	"sync"
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
)

type recordedWrite struct {
	at           time.Time
	status, note int64
}

// timedMIDIOut records the wall-clock time of every short message.
type timedMIDIOut struct {
	mutex  sync.Mutex
	writes []recordedWrite
}

func (o *timedMIDIOut) WriteShort(status, data1, data2 int64) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.writes = append(o.writes, recordedWrite{at: time.Now(), status: status, note: data1})
	return nil
}

func (o *timedMIDIOut) Close() error { return nil }

func (o *timedMIDIOut) noteOns() (list []recordedWrite) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for _, each := range o.writes {
		if each.status == noteOn {
			list = append(list, each)
		}
	}
	return
}

// newTimelineTestContext wires a registry with a real OutputDevice and Timeline. At bpm 600 a quarter note lasts 100ms.
func newTimelineTestContext(t *testing.T) (core.Context, *timedMIDIOut) {
	t.Helper()
	out := new(timedMIDIOut)
	timeline := core.NewTimeline()
	device := NewOutputDevice(1, out, 1, timeline)
	device.Start()
	t.Cleanup(timeline.Stop)
	registry := &DeviceRegistry{
		mutex:           new(sync.RWMutex),
		in:              map[int]*InputDevice{},
		out:             map[int]*OutputDevice{1: device},
		defaultInputID:  -1,
		defaultOutputID: 1,
	}
	return core.PlayContext{
		LoopControl: core.NewBeatmaster(nil, 600),
		AudioDevice: registry,
	}, out
}

func waitForNoteOns(out *timedMIDIOut, count int, timeout time.Duration) []recordedWrite {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if ons := out.noteOns(); len(ons) >= count {
			return ons
		}
		time.Sleep(5 * time.Millisecond)
	}
	return out.noteOns()
}

func checkNoteOns(t *testing.T, ons []recordedWrite, notes []int64, spacing time.Duration) {
	t.Helper()
	if len(ons) < len(notes) {
		t.Fatalf("got %d note-ons, want %d", len(ons), len(notes))
	}
	const tolerance = 40 * time.Millisecond
	for i, want := range notes {
		if ons[i].note != want {
			t.Errorf("note-on %d is %d, want %d", i, ons[i].note, want)
		}
		if i == 0 {
			continue
		}
		gap := ons[i].at.Sub(ons[i-1].at)
		if gap < spacing-tolerance || gap > spacing+tolerance {
			t.Errorf("gap before note-on %d is %s, want %s (+/- %s)", i, gap, spacing, tolerance)
		}
	}
}

func TestSequenceThroughOutputDeviceAndTimeline(t *testing.T) {
	ctx, out := newTimelineTestContext(t)
	ctx.Device().Play(core.NoCondition, core.S("C E G"), 600, time.Now())
	ons := waitForNoteOns(out, 3, 2*time.Second)
	checkNoteOns(t, ons, []int64{60, 64, 67}, 100*time.Millisecond)
}

func TestLoopThroughOutputDeviceAndTimeline(t *testing.T) {
	ctx, out := newTimelineTestContext(t)
	loop := core.NewLoop(ctx, []core.Sequenceable{core.S("C E G")})
	loop.Play(ctx, core.NoCondition, time.Now())
	t.Cleanup(func() { loop.Stop(ctx) })

	// three cycles
	ons := waitForNoteOns(out, 9, 3*time.Second)
	checkNoteOns(t, ons, []int64{60, 64, 67, 60, 64, 67, 60, 64, 67}, 100*time.Millisecond)
}
