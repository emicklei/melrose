// +build !udp

package transport

import (
	"fmt"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/tre"
	"github.com/rakyll/portmidi"
)

func init() {
	if core.IsDebug() {
		notify.Debugf("transport.init with PortmidiTransporter")
	}
	Factory = func() Transporter {
		if err := portmidi.Initialize(); err != nil {
			notify.Print(notify.Warningf("%v", tre.New(err, "portmidi.Initialize")))
		}
		return new(PortmidiTransporter)
	}
}

type PortmidiTransporter struct {
}

func (t *PortmidiTransporter) NewMIDIOut(id int) (MIDIOut, error) {
	return portmidi.NewOutputStream(portmidi.DeviceID(id), 1024, 0) // TODO flag
}

func (t *PortmidiTransporter) NewMIDIIn(id int) (MIDIIn, error) {
	return portmidi.NewInputStream(portmidi.DeviceID(id), 1024) // TODO flag
}

func (t *PortmidiTransporter) HasInputCapability() bool {
	return int(portmidi.DefaultInputDeviceID()) != -1
}

func (t *PortmidiTransporter) Terminate() {
	portmidi.Terminate()
}

func (t *PortmidiTransporter) DefaultOutputDeviceID() int {
	return int(portmidi.DefaultOutputDeviceID())
}

func (t *PortmidiTransporter) NewMIDIListener(in MIDIIn) MIDIListener {
	return newListener(in.(*portmidi.Stream))
}

func (t *PortmidiTransporter) PrintInfo(inID, outID int) {
	notify.PrintHighlighted("usage:")
	fmt.Println(":m echo                --- toggle printing the notes that are send")
	fmt.Println(":m in      <device-id> --- change the default MIDI input  device id")
	fmt.Println(":m out     <device-id> --- change the default MIDI output device id")
	fmt.Println()

	notify.PrintHighlighted("available:")
	var midiDeviceInfo *portmidi.DeviceInfo
	for i := 0; i < portmidi.CountDevices(); i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		fmt.Printf("[midi] device %d: ", i)
		usage := "output"
		if midiDeviceInfo.IsInputAvailable {
			usage = "input"
		}
		oc := "open"
		if !midiDeviceInfo.IsOpened {
			oc = "closed"
		}
		fmt.Print("\"", midiDeviceInfo.Interface, "/", midiDeviceInfo.Name, "\"",
			", is ", oc, " for ", usage)
		fmt.Println()
	}
	fmt.Println()

	notify.PrintHighlighted("current:")

	midiDeviceInfo = portmidi.Info(portmidi.DeviceID(inID))
	fmt.Printf("[midi] device  %d = default  input, %s/%s\n", inID, midiDeviceInfo.Interface, midiDeviceInfo.Name)

	midiDeviceInfo = portmidi.Info(portmidi.DeviceID(outID))
	fmt.Printf("[midi] device  %d = default output, %s/%s\n", outID, midiDeviceInfo.Interface, midiDeviceInfo.Name)
}
