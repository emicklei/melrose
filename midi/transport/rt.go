//go:build !wasm
// +build !wasm

package transport

import (
	"github.com/emicklei/melrose/notify"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

func init() { Initializer = rtInitialize }

func rtInitialize() {
	if notify.IsDebug() {
		notify.Debugf("transport.init: use RtmidiTransporter")
	}
	Factory = func() Transporter {
		return RtmidiTransporter{}
	}
}

type RtmidiTransporter struct{}

func (t RtmidiTransporter) HasInputCapability() bool {
	return true
}

func (t RtmidiTransporter) DefaultOutputDeviceID() int {
	return 0
}
func (t RtmidiTransporter) DefaultInputDeviceID() int {
	return 0
}

func (t RtmidiTransporter) NewMIDIOut(id int) (MIDIOut, error) {
	out, err := rtmidi.NewMIDIOutDefault()
	if err != nil {
		return nil, err
	}
	err = out.OpenPort(id, "")
	if err != nil {
		return nil, err
	}
	return RtmidiOut{out: out, port: id}, nil
}
func (t RtmidiTransporter) NewMIDIIn(id int) (MIDIIn, error) {
	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		return nil, err
	}
	err = in.OpenPort(id, "")
	if err != nil {
		return nil, err
	}
	// Ignore sysex, timing, or active sensing messages.
	in.IgnoreTypes(true, true, true)
	return RtmidiIn{in: in, port: id}, nil
}
func (t RtmidiTransporter) NewMIDIListener(in MIDIIn) MIDIListener {
	return newRtListener(in.(RtmidiIn).in)
}

type RtmidiOut struct {
	out  rtmidi.MIDIOut
	port int
}

func (o RtmidiOut) WriteShort(status int64, data1 int64, data2 int64) error {
	return o.out.SendMessage([]byte{byte(status & 0xFF), byte(data1 & 0xFF), byte(data2 & 0xFF)})
}
func (o RtmidiOut) Close() error {
	if notify.IsDebug() {
		name, _ := o.out.PortName(o.port)
		notify.Debugf("transport.RtmidiOut.Close: name=%s port=%d", name, o.port)
	}
	return o.out.Close()
}

type RtmidiIn struct {
	in   rtmidi.MIDIIn
	port int
}

func (i RtmidiIn) Close() error {
	if notify.IsDebug() {
		name, _ := i.in.PortName(i.port)
		notify.Debugf("transport.RtmidiIn.Close: name=%s port=%d", name, i.port)
	}
	return i.in.Close()
}

type RtListener struct {
	*mListener
	midiIn rtmidi.MIDIIn
}

func newRtListener(in rtmidi.MIDIIn) *RtListener {
	lis := &RtListener{
		midiIn:    in,
		mListener: newMListener(),
	}
	// since l.midiIn.SetCallback is blocking on success, there is no meaningful way to get an error
	// and set the callback non blocking
	go func(l *RtListener) {
		if err := l.midiIn.SetCallback(l.handleRtEvent); err != nil {
			notify.Warnf("failed to set listener callback")
		}
	}(lis)
	return lis
}

func (l *RtListener) Start() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.listening {
		return
	}
	l.listening = true

}

func (l *RtListener) handleRtEvent(m rtmidi.MIDIIn, data []byte, delta float64) {
	if !l.listening {
		// consume event so the queue does not fill up

		// if notify.IsDebug() {
		// 	notify.Debugf("handle.rt.event:%v,%v\n", data, delta)
		// }
		return
	}
	if len(data) != 3 {
		return
	}
	status := int16(data[0])
	nr := int(data[1])
	data2 := int(data[2])
	l.HandleMIDIMessage(status, nr, data2)
}

func (l *RtListener) Stop() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.listening {
		if err := l.midiIn.CancelCallback(); err != nil {
			notify.Warnf("failed to cancel listener callback")
		}
	}
	l.listening = false
}
