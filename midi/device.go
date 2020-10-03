package midi

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
	"github.com/rakyll/portmidi"
)

// Midi is an melrose.AudioDevice
type Midi struct {
	enabled      bool
	isHumanizing bool
	stream       *portmidi.Stream
	deviceID     int
	echo         bool // TODO remove

	defaultOutputChannel  int
	currentOutputDeviceID int
	currentInputDeviceID  int

	timeline *core.Timeline
	//humanize
	timingModifier   TimingModifier
	velocityModifier VelocityModifier
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
	controlChange int64 = 0xB0 // 10110000 , 176
	noteAllOff    int64 = 0x78 // 01111000 , 120  (not 123 because sustain)
	sustainPedal  int64 = 0x40
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

func (m *Midi) Timeline() *core.Timeline { return m.timeline }

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
	case "humanize":
		m.isHumanizing = !m.isHumanizing
		m.setupHumanizing()
		return notify.Infof("humanizing notes enabled:%v", m.isHumanizing)
	case "echo":
		echoMIDISent = !echoMIDISent
		return notify.Infof("printing notes enabled:%v", echoMIDISent)
	case "channel":
		if len(args) != 2 {
			return notify.Warningf("missing channel number")
		}
		nr, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.Errorf("bad channel number:%v", err)
		}
		if nr < 1 || nr > 16 {
			return notify.Errorf("bad channel number; must be in [1..16]")
		}
		m.defaultOutputChannel = nr
		return nil
	case "in":
		if len(args) != 2 {
			return notify.Warningf("missing device number")
		}
		nr, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		m.currentInputDeviceID = nr
		return notify.Infof("Current input device id:%v", m.currentInputDeviceID)
	case "out":
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
	case "init":
		m.Close()
		m.init()
		m.printInfo()
		return notify.Infof("MIDI re-initialized")
	default:
		return notify.Warningf("unknown midi command: %s", args[0])
	}
}

func (m *Midi) printInfo() {
	fmt.Println("Usage:")
	fmt.Println(":m echo                --- toggle printing the notes that are send")
	fmt.Println(":m humanize            --- toggle humanizing the notes")
	fmt.Println(":m in      <device-id> --- change the current MIDI input device id")
	fmt.Println(":m out     <device-id> --- change the current MIDI output device id")
	fmt.Println(":m channel <1..16>     --- change the default MIDI output channel")
	fmt.Println(":m init                --- initialize MIDI and device list")
	fmt.Println()

	var midiDeviceInfo *portmidi.DeviceInfo
	defaultOut := portmidi.DefaultOutputDeviceID()
	defaultIn := portmidi.DefaultInputDeviceID()

	for i := 0; i < portmidi.CountDevices(); i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		fmt.Printf("[midi] device id = %d: ", i)
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
	fmt.Printf("[midi] %v = echo notes\n", m.echo)
	fmt.Printf("[midi] %v = humanizing\n", m.isHumanizing)
	fmt.Printf("[midi] %d = input  device id (default = %d)\n", m.currentInputDeviceID, defaultIn)
	fmt.Printf("[midi] %d = output device id (default = %d)\n", m.currentOutputDeviceID, defaultOut)
	fmt.Printf("[midi] %d = default output channel\n", m.defaultOutputChannel)
}

func Open() (*Midi, error) {
	m := new(Midi)
	if err := m.init(); err != nil {
		return nil, err
	}
	m.echo = false
	// for output
	m.defaultOutputChannel = DefaultChannel
	// start timeline
	m.timeline = core.NewTimeline()
	m.setupHumanizing()
	go m.timeline.Play()
	return m, nil
}

func (m *Midi) init() error {
	portmidi.Initialize()
	deviceID := portmidi.DefaultOutputDeviceID()
	if deviceID == -1 {
		return errors.New("no default output MIDI device available")
	}
	m.enabled = true
	m.changeOutputDeviceID(int(portmidi.DefaultOutputDeviceID()))
	return nil
}

func (m *Midi) setupHumanizing() {
	if m.isHumanizing {
		// TODO make this configurable
		m.velocityModifier = newVelocityOffset(-5, 5)
		m.timingModifier = newTimingOffset(-20, 20, -20, 20)
	} else {
		m.timingModifier = NoOffset{}
		m.velocityModifier = NoOffset{}
	}
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
	return notify.Infof("MIDI device output id:%d", id)
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
	m.enabled = false
}

// echo -e "\033[93mred\033[m" # Prints “red” in red.

// 93 is bright yellow
func print(arg interface{}) {
	fmt.Printf("\033[2;93m" + fmt.Sprintf("%v ", arg) + "\033[0m")
}

func info(arg interface{}) {
	fmt.Printf("\033[2;33m" + fmt.Sprintf("%v ", arg) + "\033[0m")
}
