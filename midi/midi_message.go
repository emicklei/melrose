package midi

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type Message struct {
	device  *Device
	status  int
	channel core.Valueable
	data1   core.Valueable
	data2   core.Valueable
}

func NewMessage(device *Device, status int, channel, data1, data2 core.Valueable) Message {
	return Message{device: device, status: status, channel: channel, data1: data1, data2: data2}
}

// S has the side effect that the MIDI message is send using the device of the context
func (m Message) S() core.Sequence {
	// post creation checks
	channel := core.Int(m.channel)
	data1 := core.Int(m.data1)
	data2 := core.Int(m.data2)
	if core.IsDebug() {
		notify.Debugf("midi.message: status=%d channel=%v data1=%v data2=%v", m.status, channel, data1, data2)
	}
	if err := m.device.SendRaw(m.status, channel, data1, data2); err != nil {
		notify.Console.Errorf("failed to send MIDI message(status=%d,channel=%v,data1=%v,data2=%v) error:%v",
			m.status, m.channel, m.data1, m.data2, err)
	}
	return core.EmptySequence
}

func (m Message) Storex() string {
	return fmt.Sprintf("midi_send(%d,%v,%v,%v)", m.status, core.Storex(m.channel), core.Storex(m.data1), core.Storex(m.data2))
}

// Evaluate implements core.Evaluateable
// perform the message send
func (m Message) Evaluate() error {
	m.S()
	return nil
}
