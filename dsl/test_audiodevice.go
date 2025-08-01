package dsl

import (
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

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
func (t testAudioDevice) Report()                                              {}
func (t testAudioDevice) Close() error                                         { return nil }
func (t testAudioDevice) ListDevices() (list []core.DeviceDescriptor)          { return }
