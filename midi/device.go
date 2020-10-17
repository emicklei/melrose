package midi

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/notify"
	"github.com/rakyll/portmidi"
)

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

// Device is an melrose.AudioDevice
type Device struct {
	outputStream *portmidi.Stream
	inputStream  *portmidi.Stream
	echo         bool

	defaultOutputChannel  int
	currentOutputDeviceID int
	currentInputDeviceID  int

	timeline *core.Timeline
	listener *listener
}

func (m *Device) IO() (inputDeviceID, outputDeviceID int) {
	return m.currentInputDeviceID, m.currentOutputDeviceID
}

func (m *Device) Timeline() *core.Timeline { return m.timeline }

// SetEchoNotes is part of melrose.AudioDevice
func (m *Device) SetEchoNotes(echo bool) {
	m.echo = echo
}

func (m *Device) Reset() {
	m.timeline.Reset()
	if m.outputStream != nil {
		// send note off all to all channels for current device
		for c := 1; c <= 16; c++ {
			if err := m.outputStream.WriteShort(controlChange|int64(c-1), noteAllOff, 0); err != nil {
				fmt.Println("portmidi write error:", err)
			}
		}
	}
}

// Command is part of melrose.AudioDevice
func (m *Device) Command(args []string) notify.Message {
	if len(args) == 0 {
		m.printInfo()
		return nil
	}
	switch args[0] {
	case "echo":
		m.echo = !m.echo
		return notify.Infof("printing notes enabled:%v", m.echo)
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
		if err := m.ChangeInputDeviceID(nr); err != nil {
			return notify.Error(err)
		}
		return notify.Infof("Current input device id:%v", m.currentInputDeviceID)
	case "out":
		if len(args) != 2 {
			return notify.Warningf("missing device number")
		}
		nr, err := strconv.Atoi(args[1])
		if err != nil {
			return notify.Errorf("bad device number:%v", err)
		}
		if err := m.ChangeOutputDeviceID(nr); err != nil {
			return notify.Error(err)
		}
		return notify.Infof("Current output device id:%v", m.currentOutputDeviceID)
	default:
		return notify.Warningf("unknown midi command: %s", args[0])
	}
}

func (m *Device) printInfo() {
	fmt.Println("Usage:")
	fmt.Println(":m echo                --- toggle printing the notes that are send")
	fmt.Println(":m in      <device-id> --- change the current MIDI input device id")
	fmt.Println(":m out     <device-id> --- change the current MIDI output device id")
	fmt.Println(":m channel <1..16>     --- change the default MIDI output channel")
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
	fmt.Printf("[midi] %d = input  device id (default = %d)\n", m.currentInputDeviceID, defaultIn)
	fmt.Printf("[midi] %d = output device id (default = %d)\n", m.currentOutputDeviceID, defaultOut)
	fmt.Printf("[midi] %d = default output channel\n", m.defaultOutputChannel)
}

func Open(ctx core.Context) (*Device, error) {
	m := new(Device)
	m.currentInputDeviceID = -1
	m.currentOutputDeviceID = -1
	m.timeline = core.NewTimeline()
	m.listener = newListener(ctx)
	if err := m.init(); err != nil {
		m.Close()
		return nil, err
	}
	m.echo = false
	m.defaultOutputChannel = DefaultChannel
	// continuously send output
	go m.timeline.Play()
	return m, nil
}

func (m *Device) init() error {
	portmidi.Initialize()
	// on startup, only setup output
	outputID := portmidi.DefaultOutputDeviceID()
	if outputID == -1 {
		return errors.New("no default output MIDI device available")
	}
	if err := m.ChangeOutputDeviceID(int(outputID)); err != nil {
		return err
	}
	return nil
}

// ChangeInputDeviceID close the existing and opens a new stream. id can be -1
func (m *Device) ChangeInputDeviceID(id int) error {
	if m.currentInputDeviceID == id {
		// check stream
		if m.inputStream != nil {
			return nil
		}
	}
	if id == -1 {
		if m.inputStream != nil {
			// stop listener
			m.listener.stop()
			_ = m.inputStream.Close()
		}
		return nil
	}
	// open new
	in, err := portmidi.NewInputStream(portmidi.DeviceID(id), 1024) // TODO flag
	if err != nil {
		return err
	}
	if m.inputStream != nil {
		// stop listener
		m.listener.stop()
		_ = m.inputStream.Close()
	}
	m.inputStream = in
	m.currentInputDeviceID = id
	// start listener with new stream
	m.listener.stream = m.inputStream
	go m.listener.listen()
	return nil
}

// ChangeOutputDeviceID close the existing and opens a new stream. id can be -1
func (m *Device) ChangeOutputDeviceID(id int) error {
	if m.currentOutputDeviceID == id {
		// check stream
		if m.outputStream != nil {
			return nil
		}
	}
	if id == -1 {
		if m.outputStream != nil {
			_ = m.outputStream.Close()
		}
		return nil
	}
	// open new
	out, err := portmidi.NewOutputStream(portmidi.DeviceID(id), 1024, 0) // TODO flag
	if err != nil {
		return err
	}
	if m.outputStream != nil {
		_ = m.outputStream.Close()
	}
	m.outputStream = out
	m.currentOutputDeviceID = id
	return nil
}

// Close is part of melrose.AudioDevice
func (m *Device) Close() {
	if m.timeline != nil {
		m.timeline.Reset()
	}
	if m.outputStream != nil {
		m.outputStream.Abort()
		m.outputStream.Close()
	}
	if m.inputStream != nil {
		m.inputStream.Abort()
		m.inputStream.Close()
	}
	m.listener.stop()
	portmidi.Terminate()
}
