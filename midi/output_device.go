package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type OutputDevice struct {
	id             int
	stream         MIDIOut
	defaultChannel int

	echo     bool
	timeline *core.Timeline
}

func NewOutputDevice(id int, out MIDIOut, ch int) *OutputDevice {
	return &OutputDevice{
		id:             id,
		stream:         out,
		defaultChannel: ch,
		echo:           false,
		timeline:       core.NewTimeline(),
	}
}

func (d *OutputDevice) Start() {
	go d.timeline.Play()
}

func (d *OutputDevice) Reset() {
	d.timeline.Reset()
	if core.IsDebug() {
		notify.Debugf("device.%d: sending Note OFF to all 16 channels", d.id)
	}
	if d.stream != nil {
		// send note off all to all channels for current device
		for c := 1; c <= 16; c++ {
			if err := d.stream.WriteShort(controlChange|int64(c-1), noteAllOff, 0); err != nil {
				notify.Console.Errorf("device.%d: portmidi write error:%v", d.id, err)
			}
		}
	}
}

func (d *OutputDevice) handledPedalChange(condition core.Condition, channel int, timeline *core.Timeline, moment time.Time, group []core.Note) bool {
	if len(group) == 0 || len(group) > 1 {
		return false
	}
	note := group[0]
	switch {
	case note.IsPedalUp():
		timeline.Schedule(pedalEvent{
			goingDown:  false,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
		return true
	case note.IsPedalUpDown():
		timeline.Schedule(pedalEvent{
			goingDown:  false,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
		timeline.Schedule(pedalEvent{
			goingDown:  true,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
		return true
	case note.IsPedalDown():
		timeline.Schedule(pedalEvent{
			goingDown:  true,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
	}
	return false
}
