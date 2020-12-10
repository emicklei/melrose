// +build !windows

package transport

import (
	"fmt"

	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/tre"
	"github.com/rakyll/portmidi"
)

func init() {
	//log.Println("port init")
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
func (t *PortmidiTransporter) Start() {}
func (t *PortmidiTransporter) Stop()  {}

func (t *PortmidiTransporter) PrintInfo(inID, outID int) {
	fmt.Println("\033[1;33mUsage:\033[0m")
	fmt.Println(":m echo                --- toggle printing the notes that are send")
	fmt.Println(":m in      <device-id> --- change the default MIDI input  device id")
	fmt.Println(":m out     <device-id> --- change the default MIDI output device id")
	fmt.Println()

	fmt.Println("\033[1;33mAvailable:\033[0m")
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

	fmt.Println("\033[1;33mCurrent:\033[0m")

	midiDeviceInfo = portmidi.Info(portmidi.DeviceID(inID))
	fmt.Printf("[midi] device  %d = default  input, %s/%s\n", inID, midiDeviceInfo.Interface, midiDeviceInfo.Name)

	midiDeviceInfo = portmidi.Info(portmidi.DeviceID(outID))
	fmt.Printf("[midi] device  %d = default output, %s/%s\n", outID, midiDeviceInfo.Interface, midiDeviceInfo.Name)
}
