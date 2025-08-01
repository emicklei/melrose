package core

import (
	"time"

	"github.com/emicklei/melrose/notify"
)

var _ AudioDevice = (*AudioDeviceMock)(nil)

type AudioDeviceMock struct {
}

// Command implements the AudioDevice interface.
func (m *AudioDeviceMock) Command(cmd []string) notify.Message {
	return nil
}

// Close implements the AudioDevice interface.
func (m *AudioDeviceMock) Close() error {
	return nil
}

// DefaultDeviceIDs implements the AudioDevice interface.
func (m *AudioDeviceMock) DefaultDeviceIDs() (int, int) {
	return 0, 0
}

// HandleSetting implements the AudioDevice interface.
func (m *AudioDeviceMock) HandleSetting(setting string, value []any) error {
	return nil
}

// HasInputCapability implements the AudioDevice interface.
func (m *AudioDeviceMock) HasInputCapability() bool {
	return false
}

// Play implements the AudioDevice interface.
func (m *AudioDeviceMock) Play(condition Condition, seq Sequenceable, bpm float64, beginAt time.Time) time.Time {
	return beginAt
}

// Listen implements the AudioDevice interface.
func (m *AudioDeviceMock) Listen(deviceID int, who NoteListener, isStart bool) {
	// no-op
}

// OnKey implements the AudioDevice interface.
func (m *AudioDeviceMock) OnKey(ctx Context, deviceID int, channel int, note Note, fun HasValue) error {
	return nil
}

// Schedule implements the AudioDevice interface.
func (m *AudioDeviceMock) Schedule(event TimelineEvent, beginAt time.Time) {
	// no-op
}

// ListDevices implements the AudioDevice interface.
func (m *AudioDeviceMock) ListDevices() []DeviceDescriptor {
	return nil
}

// Reset implements the AudioDevice interface.
func (m *AudioDeviceMock) Reset() {
	// no-op
}

// Report implements the AudioDevice interface.
func (m *AudioDeviceMock) Report() {
	// no-op
}
