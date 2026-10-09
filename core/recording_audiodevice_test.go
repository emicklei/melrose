package core

import (
	"testing"
	"time"
)

func recordingContext(device AudioDevice) PlayContext {
	return PlayContext{LoopControl: &TestLooper{Biab: 4}, AudioDevice: device} // bpm 120
}

func TestRecordingAudioDevicePlaysSequenceInOrder(t *testing.T) {
	dev := NewRecordingAudioDevice()
	begin := time.Now()
	end := dev.Play(NoCondition, S("C E G"), 120, begin)
	if got, want := end.Sub(begin), 1500*time.Millisecond; got != want {
		t.Errorf("end=%s, want=%s", got, want)
	}
	if len(dev.Events()) != 0 {
		t.Fatal("nothing should be recorded before events are fired")
	}
	dev.Drain()
	want := []RecordedEvent{
		{At: 0, Channel: 1, MIDI: 60, Velocity: Normal, On: true},
		{At: 500 * time.Millisecond, Channel: 1, MIDI: 60, Velocity: Normal, On: false},
		{At: 500 * time.Millisecond, Channel: 1, MIDI: 64, Velocity: Normal, On: true},
		{At: 1000 * time.Millisecond, Channel: 1, MIDI: 64, Velocity: Normal, On: false},
		{At: 1000 * time.Millisecond, Channel: 1, MIDI: 67, Velocity: Normal, On: true},
		{At: 1500 * time.Millisecond, Channel: 1, MIDI: 67, Velocity: Normal, On: false},
	}
	checkRecorded(t, dev.Events(), want)
}

func TestRecordingAudioDeviceConditionSkipsNoteOn(t *testing.T) {
	dev := NewRecordingAudioDevice()
	dev.Play(func() bool { return false }, S("C"), 120, time.Now())
	dev.Drain()
	for _, each := range dev.Events() {
		if each.On {
			t.Errorf("unexpected %s", each)
		}
	}
}

func TestRecordingAudioDeviceRunUntil(t *testing.T) {
	dev := NewRecordingAudioDevice()
	dev.Play(NoCondition, S("C E G"), 120, time.Now())
	dev.RunUntil(500 * time.Millisecond)
	if got, want := len(dev.NoteOns()), 2; got != want {
		t.Errorf("note-ons=%d, want=%d", got, want)
	}
	if dev.Pending() == 0 {
		t.Error("events should be pending")
	}
}

func TestLoopWithLoopCountOnRecordingDevice(t *testing.T) {
	leader := runningLoop
	runningLoop = nil
	t.Cleanup(func() { runningLoop = leader })

	dev := NewRecordingAudioDevice()
	ctx := recordingContext(dev)
	loop := NewLoop(ctx, []Sequenceable{S("C E G")})
	loop.SetLoopCount(3)
	loop.Play(ctx, NoCondition, time.Now())
	dev.Drain()

	ons := dev.NoteOns()
	if got, want := len(ons), 9; got != want {
		t.Fatalf("note-ons=%d, want=%d", got, want)
	}
	for i, each := range ons {
		if want := time.Duration(i) * 500 * time.Millisecond; each.At != want {
			t.Errorf("note-on %d at %s, want %s", i, each.At, want)
		}
		if want := []int{60, 64, 67}[i%3]; each.MIDI != want {
			t.Errorf("note-on %d is %d, want %d", i, each.MIDI, want)
		}
	}
	if loop.IsRunning() {
		t.Error("loop should have ended after its loop count")
	}
}

func checkRecorded(t *testing.T, got, want []RecordedEvent) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d events, want %d:\n%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("event %d: got [%s], want [%s]", i, got[i], want[i])
		}
	}
}
