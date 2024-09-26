package dsl

import (
	"strings"
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func newTestEvaluator() *Evaluator {
	return NewEvaluator(testContext())
}

func testContext() core.Context {
	return core.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     core.NoLooper,
		AudioDevice:     testAudioDevice{},
	}
}

var _ core.AudioDevice = (*testAudioDevice)(nil)

type testAudioDevice struct{}

func (t testAudioDevice) Command(args []string) notify.Message { return nil }
func (t testAudioDevice) DefaultDeviceIDs() (int, int)         { return 1, 1 }
func (t testAudioDevice) Play(condition core.Condition, seq core.Sequenceable, bpm float64, beginAt time.Time) (endingAt time.Time) {
	return time.Now()
}
func (t testAudioDevice) HandleSetting(name string, values []interface{}) error        { return nil }
func (t testAudioDevice) HasInputCapability() bool                                     { return true }
func (t testAudioDevice) Listen(deviceID int, who core.NoteListener, startOrStop bool) {}
func (t testAudioDevice) OnKey(ctx core.Context, deviceID int, channel int, note core.Note, fun core.HasValue) error {
	return nil
}
func (t testAudioDevice) Schedule(event core.TimelineEvent, beginAt time.Time) {}
func (t testAudioDevice) Reset()                                               {}
func (t testAudioDevice) Close() error                                         { return nil }                                  { return nil }


func checkError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func eval(t *testing.T, expression string) interface{} {
	t.Helper()
	lp := new(core.TestLooper)
	lp.SetBIAB(4)
	ctx := core.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     lp,
	}
	r, err := NewEvaluator(ctx).EvaluateProgram(expression)
	checkError(t, err)
	return r
}

func mustError(t *testing.T, expression string, substring string) {
	t.Helper()
	lp := new(core.TestLooper)
	lp.SetBIAB(4)
	ctx := core.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     lp,
	}
	_, err := NewEvaluator(ctx).EvaluateExpression(expression)
	if err == nil {
		t.Fatal("error expected")
	}
	if !strings.Contains(err.Error(), substring) {
		t.Fatalf("error message should contain [%s] but it [%s]", substring, err.Error())
	}
}

func checkStorex(t *testing.T, r interface{}, storex string) {
	t.Helper()
	if s, ok := r.(core.Storable); ok {
		if got, want := s.Storex(), storex; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	} else {
		t.Errorf("result is not Storable : [%v:%T]", r, r)
	}
}
