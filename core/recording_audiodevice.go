package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

var _ AudioDevice = (*RecordingAudioDevice)(nil)

// maxDrainEvents guards Drain against loops that never end.
const maxDrainEvents = 100000

// RecordedEvent is a MIDI note on/off captured by a RecordingAudioDevice.
type RecordedEvent struct {
	At       time.Duration // since the origin of the device
	Channel  int
	MIDI     int
	Velocity int
	On       bool
}

func (e RecordedEvent) String() string {
	state := "off"
	if e.On {
		state = "on"
	}
	return fmt.Sprintf("%s %s ch=%d midi=%d vel=%d", e.At, state, e.Channel, e.MIDI, e.Velocity)
}

type recordingScheduled struct {
	event TimelineEvent
	when  time.Time
}

// RecordingAudioDevice is an AudioDevice for tests. It runs on virtual time:
// scheduled events are only fired by RunUntil or Drain and nothing ever sleeps.
// Fired note events are recorded with their time relative to the first begin time it received.
// Pedal changes are not recorded.
type RecordingAudioDevice struct {
	mutex          sync.Mutex
	origin         time.Time
	hasOrigin      bool
	queue          []recordingScheduled // ordered by time, FIFO for equal times
	recorded       []RecordedEvent
	defaultChannel int
}

func NewRecordingAudioDevice() *RecordingAudioDevice {
	return &RecordingAudioDevice{defaultChannel: 1}
}

// Events returns a copy of all fired note on/off events in the order they were fired.
func (d *RecordingAudioDevice) Events() []RecordedEvent {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return append([]RecordedEvent(nil), d.recorded...)
}

// NoteOns returns only the note-on events.
func (d *RecordingAudioDevice) NoteOns() (list []RecordedEvent) {
	for _, each := range d.Events() {
		if each.On {
			list = append(list, each)
		}
	}
	return
}

// Pending returns the number of scheduled events that have not been fired yet.
func (d *RecordingAudioDevice) Pending() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return len(d.queue)
}

// RunUntil fires all scheduled events up to and including the given offset from the origin.
func (d *RecordingAudioDevice) RunUntil(offset time.Duration) {
	for {
		d.mutex.Lock()
		if len(d.queue) == 0 || d.queue[0].when.Sub(d.origin) > offset {
			d.mutex.Unlock()
			return
		}
		head := d.queue[0]
		d.queue = d.queue[1:]
		d.mutex.Unlock()
		// without lock because handling may schedule new events
		head.event.Handle(nil, head.when)
	}
}

// Drain fires all scheduled events, including the ones that are scheduled while draining.
// It panics if there seems no end, e.g. for a Loop without a loop count.
func (d *RecordingAudioDevice) Drain() {
	for fired := 0; d.Pending() > 0; fired++ {
		if fired >= maxDrainEvents {
			panic("RecordingAudioDevice.Drain: too many events, use RunUntil or set a loop count")
		}
		d.mutex.Lock()
		head := d.queue[0]
		d.queue = d.queue[1:]
		d.mutex.Unlock()
		head.event.Handle(nil, head.when)
	}
}

func (d *RecordingAudioDevice) Play(condition Condition, seq Sequenceable, bpm float64, beginAt time.Time) time.Time {
	clock := NewPlaybackClock(beginAt, bpm)
	return d.PlayWithClock(condition, seq, &clock)
}

func (d *RecordingAudioDevice) PlayWithClock(condition Condition, seq Sequenceable, clock *PlaybackClock) time.Time {
	d.setOrigin(clock.Time())
	seq = UnValue(seq)
	if sel, ok := seq.(DeviceSelector); ok {
		seq = sel.Unwrap()
	}
	channel := d.defaultChannel
	if sel, ok := seq.(ChannelSelector); ok {
		channel = sel.Channel()
		seq = sel.Unwrap()
	}
	for _, group := range seq.S().Notes {
		if len(group) == 0 {
			continue
		}
		moment := clock.Time()
		for _, each := range group {
			if each.IsRest() || each.IsPedal() {
				continue
			}
			on := recordingNoteEvent{device: d, channel: channel, midi: each.MIDI(), velocity: each.Velocity, isOn: true, condition: condition}
			off := on
			off.isOn = false
			d.Schedule(on, moment)
			d.Schedule(off, clock.After(each).Time())
		}
		*clock = clock.AfterGroup(group)
	}
	return clock.Time()
}

// Schedule is part of AudioDevice.
func (d *RecordingAudioDevice) Schedule(event TimelineEvent, beginAt time.Time) {
	d.setOrigin(beginAt)
	d.mutex.Lock()
	defer d.mutex.Unlock()
	i := len(d.queue)
	for i > 0 && d.queue[i-1].when.After(beginAt) {
		i--
	}
	d.queue = append(d.queue, recordingScheduled{})
	copy(d.queue[i+1:], d.queue[i:])
	d.queue[i] = recordingScheduled{event: event, when: beginAt}
}

// Reset forgets all scheduled events, like the MIDI device does. Recorded events are kept.
func (d *RecordingAudioDevice) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.queue = nil
}

func (d *RecordingAudioDevice) setOrigin(t time.Time) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !d.hasOrigin {
		d.origin, d.hasOrigin = t, true
	}
}

func (d *RecordingAudioDevice) record(when time.Time, channel, midi, velocity int, on bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.recorded = append(d.recorded, RecordedEvent{At: when.Sub(d.origin), Channel: channel, MIDI: midi, Velocity: velocity, On: on})
}

type recordingNoteEvent struct {
	device    *RecordingAudioDevice
	channel   int
	midi      int
	velocity  int
	isOn      bool
	condition Condition
}

func (e recordingNoteEvent) Handle(tim *Timeline, when time.Time) {
	if e.isOn && e.condition != nil && !e.condition() {
		return
	}
	e.device.record(when, e.channel, e.midi, e.velocity, e.isOn)
}

func (e recordingNoteEvent) NoteChangesDo(block func(NoteChange)) {
	block(NewNoteChange(e.isOn, int64(e.midi), int64(e.velocity)))
}

func (d *RecordingAudioDevice) Command(args []string) notify.Message { return nil }
func (d *RecordingAudioDevice) DefaultDeviceIDs() (int, int)         { return 1, 1 }
func (d *RecordingAudioDevice) HandleSetting(name string, values []any) error {
	return nil
}
func (d *RecordingAudioDevice) HasInputCapability() bool                                { return false }
func (d *RecordingAudioDevice) Listen(deviceID int, who NoteListener, startOrStop bool) {}
func (d *RecordingAudioDevice) OnKey(ctx Context, deviceID int, channel int, note Note, fun HasValue) error {
	return nil
}
func (d *RecordingAudioDevice) ListDevices() []DeviceDescriptor { return nil }
func (d *RecordingAudioDevice) Report()                         {}
func (d *RecordingAudioDevice) Close() error                    { return nil }
