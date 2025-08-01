package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/emicklei/melrose/notify"
)

func Test_tickerDuration(t *testing.T) {
	d := beatTickerDuration(60.0)
	if got, want := d, 1*time.Second; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	d2 := beatTickerDuration(300.0)
	if got, want := d2, 200*time.Millisecond; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	d3 := beatTickerDuration(120.0)
	if got, want := d3, 500*time.Millisecond; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBeatmaster_beatsAtNextBar(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.beats = 0
	b.biab = 3
	if got, want := b.beatsAtNextBar(), int64(0); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	b.beats = 5
	b.biab = 4
	if got, want := b.beatsAtNextBar(), int64(8); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTrackBarTiming(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.SetBPM(100.0)
	b.SetBIAB(3)
	ctx.LoopControl = b

	tr1 := NewTrack("1", 1)
	tr1.Add(NewSequenceOnTrack(On(1), MustParseSequence("c d e")))
	tr2 := NewTrack("2", 1)
	tr1.Add(NewSequenceOnTrack(On(2), MustParseSequence("c")))
	m := MultiTrack{Tracks: []HasValue{On(tr1), On(tr2)}}
	m.Play(ctx, time.Now())
	t.Log(b.schedule.entries)
	_, ok := b.schedule.entries[0]
	if !ok {

		t.Fail()
	}
	_, ok = b.schedule.entries[3]
	if !ok {
		t.Fail()
	}
}

func TestBeatmaster_Getters(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	if got, want := b.BPM(), 120.0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := b.BIAB(), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if g1, g2 := b.BeatsAndBars(); g1 != 0 || g2 != 0 {
		t.Errorf("got [%v,%v] want [%v,%v]", g1, g2, 0, 0)
	}
}

func TestBeatmaster_Setters(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.SetBPM(140.0)
	b.SetBIAB(3)
	if got, want := b.BPM(), 140.0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := b.BIAB(), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBeatmaster_StartStop(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	notify.ToggleDebug()
	b.Start()
	if !b.beating {
		t.Error("not beating")
	}
	b.Start() // should do nothing
	b.Stop()
	if b.beating {
		t.Error("still beating")
	}
	b.Stop() // should do nothing
	notify.ToggleDebug()
}

func TestBeatmaster_Reset(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.Start()
	b.Reset()
	if !b.beating {
		t.Error("not beating")
	}
}

func TestBeatmaster_SetBPM_Beating(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.Start()
	defer b.Stop()
	b.SetBPM(140.0)
	// same bpm
	b.SetBPM(140.0)
}

func TestBeatmaster_SetBIAB_Beating(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.Start()
	defer b.Stop()
	b.SetBIAB(3)
	if got, want := b.BIAB(), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	// same biab
	b.SetBIAB(3)
}

func TestBeatmaster_Plan(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.Plan(1, MustParseSequence("C"))
	if got, want := len(b.schedule.entries), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	// test debug
	notify.ToggleDebug()
	b.Plan(1, MustParseSequence("D"))
	notify.ToggleDebug()

	// test nil device
	b.context = Context(nil)
	b.Plan(1, MustParseSequence("E"))
}

func TestBeatmaster_SettingNotifier(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.Start()
	defer b.Stop()
	notified := false
	b.SettingNotifier(func(LoopController) {
		notified = true
	})
	b.SetBIAB(2)
	if !notified {
		t.Error("not notified")
	}
}

func TestBeatmaster_SettingNotifierBPM(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	notified := false
	b.SettingNotifier(func(LoopController) {
		notified = true
	})
	b.SetBPM(2)
	if !notified {
		t.Error("not notified")
	}
}

type mockDevice struct {
	played bool
}

func (d *mockDevice) DefaultDeviceIDs() (inputDeviceID, outputDeviceID int) {
	return 0, 0
}
func (d *mockDevice) Command(args []string) notify.Message {
	return nil
}
func (d *mockDevice) HandleSetting(name string, values []any) error {
	return nil
}
func (d *mockDevice) Play(condition Condition, seq Sequenceable, bpm float64, beginAt time.Time) (endingAt time.Time) {
	d.played = true
	return time.Now()
}
func (d *mockDevice) HasInputCapability() bool {
	return false
}
func (d *mockDevice) Listen(deviceID int, who NoteListener, isStart bool) {}
func (d *mockDevice) OnKey(ctx Context, deviceID int, channel int, note Note, fun HasValue) error {
	return nil
}
func (d *mockDevice) Schedule(event TimelineEvent, beginAt time.Time) {}
func (d *mockDevice) ListDevices() []DeviceDescriptor {
	return nil
}
func (d *mockDevice) Reset() {}
func (d *mockDevice) Close() error {
	return nil
}
func (d *mockDevice) Report() {}

func TestBeatmaster_PlayPlan(t *testing.T) {
	dev := new(mockDevice)
	ctx := &PlayContext{
		AudioDevice: dev,
	}
	b := NewBeatmaster(ctx, 120.0)
	ctx.LoopControl = b
	b.Plan(0, MustParseSequence("C"))
	b.Start()
	defer b.Stop()
	time.Sleep(1 * time.Second)
	if !dev.played {
		t.Error("device was not played")
	}
}

func TestBeatmaster_ResetDrain(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(&ctx, 120.0)
	go func() {
		b.bpmChanges <- 140.0
	}()
	time.Sleep(10 * time.Millisecond)
	b.Reset()
}

func TestZeroBeat(t *testing.T) {
	z := zeroBeat{}
	z.Start()
	z.Stop()
	z.Reset()
	z.SetBPM(1)
	if got, want := z.BPM(), 120.0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	z.SetBIAB(1)
	if got, want := z.BIAB(), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if g1, g2 := z.BeatsAndBars(); g1 != 0 || g2 != 0 {
		t.Errorf("got [%v,%v] want [%v,%v]", g1, g2, 0, 0)
	}
	z.Plan(1, MustParseSequence("C"))
	z.SettingNotifier(nil)
}

func TestBeatmaster_ScheduleEmpty(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.Start()
	defer b.Stop()
	time.Sleep(2 * time.Second)
	if got, want := b.beats, int64(0); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBeatmaster_ChangeBPMOnBar(t *testing.T) {
	var buf bytes.Buffer
	notify.Console.StandardOut = &buf
	notify.Console.StandardError = &buf
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.schedule.Schedule(0, func(t time.Time) {})
	b.Start()
	defer b.Stop()
	b.SetBPM(140)
	time.Sleep(2 * time.Second)
	if got, want := b.BPM(), 140.0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
