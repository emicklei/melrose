package dsl

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
)

// recordingTestEvaluator uses virtual time and a bpm of 120, so a quarter note takes 500ms.
func recordingTestEvaluator() (*Evaluator, core.Context, *core.RecordingAudioDevice) {
	dev := core.NewRecordingAudioDevice()
	ctx := core.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     &core.TestLooper{Biab: 4},
		AudioDevice:     dev,
		EnvironmentVars: new(sync.Map),
	}
	return NewEvaluator(ctx), ctx, dev
}

// noteOns formats the note-on events as "midi@milliseconds".
func noteOns(dev *core.RecordingAudioDevice) (list []string) {
	for _, each := range dev.NoteOns() {
		list = append(list, fmt.Sprintf("%d@%d", each.MIDI, each.At.Milliseconds()))
	}
	return
}

func TestPlayExpressionsRecorded(t *testing.T) {
	for _, each := range []struct {
		expr string
		want []string
	}{
		{"play(note('C'))", []string{"60@0"}},
		{"play(seq('C E G'))", []string{"60@0", "64@500", "67@1000"}},
		{"play(seq('C = G'))", []string{"60@0", "67@1000"}},
		{"play(seq('C (E G)'))", []string{"60@0", "64@500", "67@500"}},
		{"play(seq('2C 8D'))", []string{"60@0", "62@1000"}},
		{"play(seq('C'), seq('E'))", []string{"60@0", "64@500"}},
		{"play(repeat(2, seq('C D')))", []string{"60@0", "62@500", "60@1000", "62@1500"}},
		{"play(pitch(2, seq('C E')))", []string{"62@0", "66@500"}},
	} {
		t.Run(each.expr, func(t *testing.T) {
			e, _, dev := recordingTestEvaluator()
			if _, err := e.EvaluateProgram(each.expr); err != nil {
				t.Fatal(err)
			}
			dev.Drain()
			if got := noteOns(dev); !reflect.DeepEqual(got, each.want) {
				t.Errorf("got %v, want %v", got, each.want)
			}
		})
	}
}

func TestNoteOffAfterNoteOnRecorded(t *testing.T) {
	e, _, dev := recordingTestEvaluator()
	if _, err := e.EvaluateProgram("play(seq('C 2E'))"); err != nil {
		t.Fatal(err)
	}
	dev.Drain()
	var got []string
	for _, each := range dev.Events() {
		got = append(got, fmt.Sprintf("%d:%v@%d", each.MIDI, each.On, each.At.Milliseconds()))
	}
	want := []string{"60:true@0", "60:false@500", "64:true@500", "64:false@1500"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPlayLoopExpressionRecorded(t *testing.T) {
	e, _, dev := recordingTestEvaluator()
	if _, err := e.EvaluateProgram("play(loop(seq('C E G')))"); err != nil {
		t.Fatal(err)
	}
	// cycles last 1500ms; the third cycle starts at 3000ms
	dev.RunUntil(4 * time.Second)
	want := []string{
		"60@0", "64@500", "67@1000",
		"60@1500", "64@2000", "67@2500",
		"60@3000", "64@3500", "67@4000",
	}
	if got := noteOns(dev); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPlayLoopVariableRecorded(t *testing.T) {
	e, ctx, dev := recordingTestEvaluator()
	r, err := e.EvaluateProgram("bar = loop(seq('C E G'))")
	if err != nil {
		t.Fatal(err)
	}
	loop, ok := r.(*core.Loop)
	if !ok {
		t.Fatalf("got %T, want *core.Loop", r)
	}
	loop.SetLoopCount(2)
	loop.Play(ctx, core.NoCondition, time.Now())
	dev.Drain()
	want := []string{"60@0", "64@500", "67@1000", "60@1500", "64@2000", "67@2500"}
	if got := noteOns(dev); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
