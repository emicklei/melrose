package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

// Schedule exists for Loop
func (registry *DeviceRegistry) Schedule(e core.TimelineEvent, beginAt time.Time) {
	// TODO check DeviceSelector
	device, err := registry.Output(registry.defaultOutputID)
	if err != nil {
		return
	}
	device.timeline.Schedule(e, beginAt)
}

// Play schedules all the notes on the timeline beginning at a give time (now or in the future).
// Returns the end time of the last played Note.
func (registry *DeviceRegistry) Play(condition core.Condition, seq core.Sequenceable, bpm float64, beginAt time.Time) time.Time {
	if core.IsDebug() {
		notify.Debugf("midi.play: time=%s object=%s", beginAt.Format("04:05.000"), core.Storex(seq))
	}
	// unwrap if variable because we need to detect device or channel selector
	seq = core.UnValue(seq)

	// which device?
	var device *OutputDevice
	deviceID := registry.defaultOutputID
	if dev, ok := seq.(core.DeviceSelector); ok {
		deviceID = dev.DeviceID()
		seq = dev.Unwrap()
	}
	device, err := registry.Output(deviceID)
	if err != nil {
		return beginAt
	}

	return device.Play(condition, seq, bpm, beginAt)
}
