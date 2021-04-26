// +build wasm

package transport

import (
	"syscall/js"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

func init() { Initializer = wasmInitialize }

func wasmInitialize() {
	if core.IsDebug() {
		notify.Debugf("transport.init: use WASMmidiTransporter")
	}
	Factory = func() Transporter {
		return WASMmidiTransporter{}
	}
}

type WASMmidiTransporter struct{}

func (t WASMmidiTransporter) HasInputCapability() bool {
	return true
}
func (t WASMmidiTransporter) PrintInfo(inID, outID int) {

}
func (t WASMmidiTransporter) DefaultOutputDeviceID() int {
	return 0
}
func (t WASMmidiTransporter) DefaultInputDeviceID() int {
	return 0
}
func (t WASMmidiTransporter) NewMIDIOut(id int) (MIDIOut, error) {
	return WASMMidiOut{id: id}, nil

}
func (t WASMmidiTransporter) NewMIDIIn(id int) (MIDIIn, error) {
	return WASMMidiIn{id: id}, nil
}
func (t WASMmidiTransporter) NewMIDIListener(m MIDIIn) MIDIListener {
	return &WASMListener{
		mListener: newMListener(),
	}
}

type WASMMidiOut struct {
	id int
}

func (m WASMMidiOut) WriteShort(status int64, data1 int64, data2 int64) error {
	// MIDI_send(deviceID, status, pitch, velocity)
	js.Global().Call("melrose_send", m.id, uint8(status), uint8(data1), uint8(data2))
	return nil
}
func (m WASMMidiOut) Close() error {
	return nil
}

type WASMMidiIn struct {
	id int
}

func (m WASMMidiIn) Close() error {
	return nil
}

type WASMListener struct {
	*mListener
}

func (l *WASMListener) Start() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.listening {
		return
	}
	l.listening = true
}

func (l *WASMListener) Stop() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.listening {
		return
	}
	l.listening = false
}

func (l *WASMListener) HandleMIDIMessage(status int16, nr int, data2 int) {
	if !l.listening {
		return
	}
	l.mListener.HandleMIDIMessage(status, nr, data2)
}
