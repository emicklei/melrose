package dsl

import (
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func Validate(source string) error {
	ctx := core.PlayContext{
		VariableStorage: NewVariableStore(),
		LoopControl:     core.NoLooper,
		AudioDevice:     noAudioDevice{},
	}
	e := NewEvaluator(ctx)
	_, err := e.EvaluateProgram(source)
	return err
}

var _ core.AudioDevice = (*noAudioDevice)(nil)

type noAudioDevice struct{}

func (t noAudioDevice) Command(args []string) notify.Message { return nil }
func (t noAudioDevice) DefaultDeviceIDs() (int, int)         { return 1, 1 }
func (t noAudioDevice) Play(condition core.Condition, seq core.Sequenceable, bpm float64, beginAt time.Time) (endingAt time.Time) {
	return time.Now()
}
func (t noAudioDevice) HandleSetting(name string, values []any) error                { return nil }
func (t noAudioDevice) HasInputCapability() bool                                     { return true }
func (t noAudioDevice) Listen(deviceID int, who core.NoteListener, startOrStop bool) {}
func (t noAudioDevice) OnKey(ctx core.Context, deviceID int, channel int, note core.Note, fun core.HasValue) error {
	return nil
}
func (t noAudioDevice) Schedule(event core.TimelineEvent, beginAt time.Time) {}
func (t noAudioDevice) Reset()                                               {}
func (t noAudioDevice) Report()                                              {}
func (t noAudioDevice) Close() error                                         { return nil }
func (t noAudioDevice) ListDevices() (list []core.DeviceDescriptor)          { return }
