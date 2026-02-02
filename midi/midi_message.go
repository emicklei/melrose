package midi

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi/transport"
	"github.com/emicklei/melrose/notify"
)

// https://www.midi.org/specifications-old/item/table-1-summary-of-midi-message
const (
	noteOn        int64 = 0x90 // 10010000 , 144
	noteOff       int64 = 0x80 // 10000000 , 128
	controlChange int64 = 0xB0 // 10110000 , 176
	noteAllOff    int64 = 0x78 // 01111000 , 120  (not 123 because sustain)
	sustainPedal  int64 = 0x40
)

type Message struct {
	audioDevices core.AudioDevice
	status       int
	deviceID     core.HasValue
	channel      core.HasValue
	data1        core.HasValue
	data2        core.HasValue
}

func NewMessage(audioDevices core.AudioDevice, id core.HasValue, status int, channel, data1, data2 core.HasValue) Message {
	return Message{audioDevices: audioDevices, deviceID: id, status: status, channel: channel, data1: data1, data2: data2}
}

// S has the side effect that the MIDI message is send using the device of the context
func (m Message) S() core.Sequence {
	// post creation checks
	deviceID := core.Int(m.deviceID)
	channel := core.Int(m.channel)
	data1 := core.Int(m.data1)
	data2 := core.Int(m.data2)
	if notify.IsDebug() {
		notify.Debugf("midi.message: device=%d, status=%d channel=%v data1=%v data2=%v", deviceID, m.status, channel, data1, data2)
	}
	devices := m.audioDevices.(*DeviceRegistry)
	out, err := devices.Output(deviceID)
	if err != nil {
		notify.Console.Errorf("failed to send MIDI message(device=%d,status=%d,channel=%v,data1=%v,data2=%v) error:%v",
			deviceID, m.status, m.channel, m.data1, m.data2, err)
		return core.EmptySequence
	}

	if err := sendRaw(m.status, channel, data1, data2, out.stream); err != nil {
		notify.Console.Errorf("failed to send MIDI message(device=%d,status=%d,channel=%v,data1=%v,data2=%v) error:%v",
			deviceID, m.status, m.channel, m.data1, m.data2, err)
	}
	return core.EmptySequence
}

func (m Message) Storex() string {
	return fmt.Sprintf("midi_send(%v,%d,%v,%v,%v)", core.Storex(m.deviceID), m.status, core.Storex(m.channel), core.Storex(m.data1), core.Storex(m.data2))
}

// Evaluate implements core.Evaluatable
// perform the message send
func (m Message) Evaluate(ctx core.Context) error {
	m.S()
	return nil
}

// if channel < then do not include that information
func sendRaw(status, channel, data1, data2 int, out transport.MIDIOut) error {
	if channel < 1 {
		return out.WriteShort(int64(status), int64(data1), int64(data2))
	}
	if channel > 16 {
		return fmt.Errorf("invalid MIDI channel:%d", channel)
	}
	return out.WriteShort(int64(status|(channel-1)), int64(data1), int64(data2))
}
