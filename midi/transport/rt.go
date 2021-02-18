package transport

import (
	"fmt"
	"log"
	"sync"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"gitlab.com/gomidi/rtmididrv/imported/rtmidi"
)

func init() { Initializer = rtInitialize }

func rtInitialize() {
	if core.IsDebug() {
		notify.Debugf("transport.init: use RtmidiTransporter")
	}
	Factory = func() Transporter {
		return &RtmidiTransporter{}
	}
}

type RtmidiTransporter struct {
	mutex sync.RWMutex
}

func (t *RtmidiTransporter) HasInputCapability() bool {
	return true
}
func (t *RtmidiTransporter) PrintInfo(inID, outID int) {
	notify.PrintHighlighted("usage:")
	fmt.Println(":m echo                               --- toggle printing the notes that are send")
	fmt.Println(":m in      <device-id>                --- change the default MIDI input  device id")
	fmt.Println(":m out     <device-id>                --- change the default MIDI output device id")
	fmt.Println(":m channel <device-id> <midi-channel> --- change the default MIDI channel for an output device id")
	fmt.Println()

	notify.PrintHighlighted("available:")

	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		log.Fatalln("can't open default MIDI in: ", err)
	}
	defer in.Close()
	ports, err := in.PortCount()
	if err != nil {
		log.Fatalln("can't get number of in ports: ", err.Error())
	}
	for i := 0; i < ports; i++ {
		name, err := in.PortName(i)
		if err != nil {
			name = ""
		}
		fmt.Println(i, name)
	}
	{
		// Outs
		out, err := rtmidi.NewMIDIOutDefault()
		if err != nil {
			log.Fatalln("can't open default MIDI out: ", err)
		}
		defer out.Close()
		ports, err := out.PortCount()
		if err != nil {
			log.Fatalln("can't get number of out ports: ", err.Error())
		}

		for i := 0; i < ports; i++ {
			name, err := out.PortName(i)
			if err != nil {
				name = ""
			}
			fmt.Println(i, name)
		}
	}
	fmt.Println()
}
func (t *RtmidiTransporter) DefaultOutputDeviceID() int {
	return 0
}
func (t *RtmidiTransporter) DefaultInputDeviceID() int {
	return 0
}

func (t *RtmidiTransporter) NewMIDIOut(id int) (MIDIOut, error) {
	out, err := rtmidi.NewMIDIOutDefault()
	if err != nil {
		return nil, err
	}
	err = out.OpenPort(id, "")
	if err != nil {
		return nil, err
	}
	return RtmidiOut{out: out}, nil
}
func (t *RtmidiTransporter) NewMIDIIn(id int) (MIDIIn, error) {
	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		return nil, err
	}
	err = in.OpenPort(id, "")
	if err != nil {
		return nil, err
	}
	return RtmidiIn{in: in}, nil
}
func (t *RtmidiTransporter) Terminate() {
	// noop
}
func (t *RtmidiTransporter) NewMIDIListener(MIDIIn) MIDIListener {
	return nil
}

type RtmidiOut struct {
	out rtmidi.MIDIOut
}

func (o RtmidiOut) WriteShort(status int64, data1 int64, data2 int64) error {
	return o.out.SendMessage([]byte{byte(status & 0xFF), byte(data1 & 0xFF), byte(data2 & 0xFF)})
}
func (o RtmidiOut) Close() error { return o.out.Close() }
func (o RtmidiOut) Abort() error { return nil }

type RtmidiIn struct {
	in rtmidi.MIDIIn
}

func (i RtmidiIn) Close() error { return i.in.Close() }
