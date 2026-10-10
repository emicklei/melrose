package api

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl"
)

// newRecordingService uses virtual time and a bpm of 120, so a quarter note takes 500ms.
func newRecordingService() (*ServiceImpl, *core.RecordingAudioDevice) {
	dev := core.NewRecordingAudioDevice()
	ctx := core.PlayContext{
		VariableStorage: dsl.NewVariableStore(),
		LoopControl:     &core.TestLooper{Biab: 4},
		AudioDevice:     dev,
		EnvironmentVars: new(sync.Map),
	}
	return NewService(ctx).(*ServiceImpl), dev
}

func noteOnTimes(dev *core.RecordingAudioDevice) (list [][2]int64) {
	for _, each := range dev.NoteOns() {
		list = append(list, [2]int64{int64(each.MIDI), each.At.Milliseconds()})
	}
	return
}

func TestCommandPlaySequenceExpression(t *testing.T) {
	svc, dev := newRecordingService()
	if _, err := svc.CommandPlay("/tmp/test.mel", 1, "play(seq('C E G'))"); err != nil {
		t.Fatal(err)
	}
	dev.Drain()
	want := [][2]int64{{60, 0}, {64, 500}, {67, 1000}}
	if got := noteOnTimes(dev); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCommandPlayBareSequence(t *testing.T) {
	svc, dev := newRecordingService()
	if _, err := svc.CommandPlay("/tmp/test.mel", 1, "seq('C E G')"); err != nil {
		t.Fatal(err)
	}
	dev.Drain()
	want := [][2]int64{{60, 0}, {64, 500}, {67, 1000}}
	if got := noteOnTimes(dev); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCommandPlayLoopVariable(t *testing.T) {
	svc, dev := newRecordingService()
	r, err := svc.CommandEvaluate("/tmp/test.mel", 1, "bar = loop(seq('C E G'))")
	if err != nil {
		t.Fatal(err)
	}
	r.(*core.Loop).SetLoopCount(2)
	if _, err := svc.CommandPlay("/tmp/test.mel", 2, "bar"); err != nil {
		t.Fatal(err)
	}
	dev.Drain()
	want := [][2]int64{{60, 0}, {64, 500}, {67, 1000}, {60, 1500}, {64, 2000}, {67, 2500}}
	if got := noteOnTimes(dev); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCommandPlayLoopExpression(t *testing.T) {
	svc, dev := newRecordingService()
	if _, err := svc.CommandPlay("/tmp/test.mel", 1, "play(loop(seq('C E G')))"); err != nil {
		t.Fatal(err)
	}
	dev.RunUntil(3 * time.Second)
	want := [][2]int64{
		{60, 0}, {64, 500}, {67, 1000},
		{60, 1500}, {64, 2000}, {67, 2500},
		{60, 3000},
	}
	if got := noteOnTimes(dev); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCommandPlayStoredPlayExpression(t *testing.T) {
	svc, dev := newRecordingService()
	if _, err := svc.CommandEvaluate("/tmp/test.mel", 1, "p = play(seq('C E G'))"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CommandPlay("/tmp/test.mel", 2, "p"); err != nil {
		t.Fatal(err)
	}
	dev.Drain()
	want := [][2]int64{{60, 0}, {64, 500}, {67, 1000}}
	if got := noteOnTimes(dev); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
