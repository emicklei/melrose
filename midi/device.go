package midi

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/emicklei/melrose/notify"
	"github.com/rakyll/portmidi"
)

// Midi is an melrose.AudioDevice
type Midi struct {
	enabled      bool
	stream       *portmidi.Stream
	mutex        *sync.Mutex
	deviceID     int
	echo         bool
	bpm          float64
	baseVelocity int

	defaultOutputChannel  int
	currentOutputDeviceID int
	currentInputDeviceID  int
}

const (
	noteOn  int = 0x90
	noteOff int = 0x80
)

var (
	DefaultVelocity = 70
	DefaultBPM      = 120.0
	DefaultChannel  = 0
)

// SetBeatsPerMinute (BPM) ; beats each the length of a quarter note per minute.
func (m *Midi) SetBeatsPerMinute(bpm float64) {
	if bpm <= 0 {
		return
	}
	m.bpm = bpm
}

// SetBaseVelocity is part of melrose.AudioDevice
func (m *Midi) SetBaseVelocity(velocity int) {
	if velocity < 0 || velocity > 127 {
		return
	}
	m.baseVelocity = velocity
}

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
	fmt.Println("MIDI: Beats per minute (bpm):", m.bpm)
	fmt.Printf("MIDI: Base velocity :%d\n", m.baseVelocity)
	fmt.Printf("MIDI: Echo notes :%v\n", m.echo)
	fmt.Println("MIDI: Default output channel:", m.defaultOutputChannel)
	var midiDeviceInfo *portmidi.DeviceInfo
	defaultOut := portmidi.DefaultOutputDeviceID()
	fmt.Println("MIDI: Default output device id:", defaultOut)
	fmt.Println("MIDI: Current output device id:", m.currentOutputDeviceID)

	defaultIn := portmidi.DefaultInputDeviceID()
	fmt.Println("MIDI: Default input device id:", defaultIn)
	fmt.Println("MIDI: Current input device id:", m.currentInputDeviceID)

	for i := 0; i < portmidi.CountDevices(); i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		fmt.Printf("MIDI: %v: ", i)
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
	m.mutex = new(sync.Mutex)
	m.enabled = false
	portmidi.Initialize()
	deviceID := portmidi.DefaultOutputDeviceID()
	if deviceID == -1 {
		return nil, errors.New("no default output MIDI device available")
	}
	m.enabled = true
	m.bpm = DefaultBPM
	m.baseVelocity = DefaultVelocity
	m.echo = true
	// for output
	m.defaultOutputChannel = DefaultChannel
	m.changeOutputDeviceID(int(portmidi.DefaultOutputDeviceID()))
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
