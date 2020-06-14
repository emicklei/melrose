package midi

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
	"github.com/rakyll/portmidi"
)

// Midi is an melrose.AudioDevice
type Midi struct {
	enabled  bool
	stream   *portmidi.Stream
	deviceID int
	echo     bool

	defaultOutputChannel  int
	currentOutputDeviceID int
	currentInputDeviceID  int

	timeline *melrose.Timeline
}

type MIDIWriter interface {
	WriteShort(int64, int64, int64) error
	Close() error
	Abort() error
}

// https://www.midi.org/specifications-old/item/table-1-summary-of-midi-message
const (
	noteOn        int64 = 0x90 // 10010000 , 144
	noteOff       int64 = 0x80 // 10000000 , 128
	controlChange int64 = 176  // 10110000 , 176
	noteAllOff    int64 = 123  // 01111011 , 123
)

var (
	DefaultChannel = 1
)

func (m *Midi) Reset() {
	m.timeline.Reset()
	if m.stream != nil {
		// send note off all to all channels for current device
		for c := 1; c <= 16; c++ {
			if err := m.stream.WriteShort(controlChange|int64(c-1), noteAllOff, 0); err != nil {
				fmt.Println("portmidi write error:", err)
			}
		}
	}
}

func (m *Midi) Timeline() *melrose.Timeline { return m.timeline }

// SetEchoNotes is part of melrose.AudioDevice
func (m *Midi) SetEchoNotes(echo bool) {
	m.echo = echo
}

// Command is part of melrose.AudioDevice
func (m *Midi) Command(args []string) notify.Message {
	if len(args) == 0 {
		m.printInfo()
		return nil
	}
	switch args[0] {
	case "input":
		if len(args) != 2 {
			return notify.Warningf("missing device number")
		}
		nr, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		m.currentInputDeviceID = nr
		return notify.Infof("Current input device id:%v", m.currentInputDeviceID)
	case "output":
		if len(args) != 2 {
			return notify.Warningf("missing device number")
		}
		nr, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		if err := m.changeOutputDeviceID(nr); err != nil {
			return err
		}
		return notify.Infof("Current output device id:%v", m.currentOutputDeviceID)
	default:
		return notify.Warningf("unknown midi command: %s", args[0])
	}
}

func (m *Midi) printInfo() {
	fmt.Println("Usage:")
	fmt.Println(":m input  <device-id> --- change the current MIDI input device id")
	fmt.Println(":m output <device-id> --- change the current MIDI output device id")
	fmt.Println()
	fmt.Printf("[midi] echo notes: %v\n", m.echo)
	fmt.Println("[midi] default output channel:", m.defaultOutputChannel)
	var midiDeviceInfo *portmidi.DeviceInfo
	defaultOut := portmidi.DefaultOutputDeviceID()
	fmt.Println("[midi] default output device id:", defaultOut)
	fmt.Println("[midi] current output device id:", m.currentOutputDeviceID)

	defaultIn := portmidi.DefaultInputDeviceID()
	fmt.Println("[midi] default input device id:", defaultIn)
	fmt.Println("[midi] current input device id:", m.currentInputDeviceID)

	for i := 0; i < portmidi.CountDevices(); i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		fmt.Printf("[midi] device id %d: ", i)
		usage := "output"
		if midiDeviceInfo.IsInputAvailable {
			usage = "input"
		}
		fmt.Print("\"", midiDeviceInfo.Interface, "/", midiDeviceInfo.Name, "\"",
			", open=", midiDeviceInfo.IsOpened,
			", usage=", usage)
		fmt.Println()
	}
}

func Open() (*Midi, error) {
	m := new(Midi)
	portmidi.Initialize()
	deviceID := portmidi.DefaultOutputDeviceID()
	if deviceID == -1 {
		return nil, errors.New("no default output MIDI device available")
	}
	m.enabled = true
	m.echo = false
	// for output
	m.defaultOutputChannel = DefaultChannel
	m.changeOutputDeviceID(int(portmidi.DefaultOutputDeviceID()))

	// start timeline
	m.timeline = melrose.NewTimeline()
	go m.timeline.Play()

	return m, nil
}

func (m *Midi) changeOutputDeviceID(id int) notify.Message {
	if !m.enabled {
		return notify.Warningf("MIDI is not enabled")
	}
	if m.currentOutputDeviceID == id {
		// check stream
		if m.stream != nil {
			return nil
		}
	}
	// open new
	out, err := portmidi.NewOutputStream(portmidi.DeviceID(id), 1024, 0)
	if err != nil {
		return notify.Error(err)
	}
	if m.stream != nil {
		// close old stream
		m.stream.Close()
	}
	m.stream = out
	m.currentOutputDeviceID = id
	return nil
}

// Close is part of melrose.AudioDevice
func (m *Midi) Close() {
	if m.timeline != nil {
		m.timeline.Reset()
	}
	if m.enabled {
		m.stream.Abort()
		m.stream.Close()
	}
	portmidi.Terminate()
}

// 94 is bright cyan
func print(arg interface{}) {
	fmt.Printf("\033[2;96m" + fmt.Sprintf("%v ", arg) + "\033[0m")
}

func info(arg interface{}) {
	fmt.Printf("\033[2;33m" + fmt.Sprintf("%v ", arg) + "\033[0m")
}
