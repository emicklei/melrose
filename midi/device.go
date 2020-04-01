package midi

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rakyll/portmidi"
)

type Midi struct {
	enabled      bool
	stream       *portmidi.Stream
	mutex        *sync.Mutex
	deviceID     int
	echo         bool
	bpm          float64
	baseVelocity int
}

const (
	noteOn          int = 0x90
	noteOff         int = 0x80
	DefaultVelocity     = 80
	DefaultBPM          = 120
)

// BeatsPerMinute (BPM) ; beats each the length of a quarter note per minute.
func (m *Midi) SetBeatsPerMinute(bpm float64) {
	if bpm <= 0 {
		return
	}
	m.bpm = bpm
}

func (m *Midi) BeatsPerMinute() float64 {
	return m.bpm
}

func (m *Midi) PrintInfo() {
	fmt.Println("[midi] BPM:", m.bpm)
	var midiDeviceInfo *portmidi.DeviceInfo
	defaultOut := portmidi.DefaultOutputDeviceID()
	fmt.Println("[midi] default output device id:", defaultOut)
	for i := 0; i < portmidi.CountDevices(); i++ {
		midiDeviceInfo = portmidi.Info(portmidi.DeviceID(i)) // returns info about a MIDI device
		fmt.Printf("[midi] %v: ", i)
		fmt.Print("\"", midiDeviceInfo.Interface, "/", midiDeviceInfo.Name, "\"")
		fmt.Println()
	}
}

func Open() (*Midi, error) {
	m := new(Midi)
	m.enabled = false
	portmidi.Initialize()
	deviceID := portmidi.DefaultOutputDeviceID()
	if deviceID == -1 {
		return nil, errors.New("no default output device available")
	}
	out, err := portmidi.NewOutputStream(deviceID, 1024, 0)
	if err != nil {
		return nil, err
	}
	m.enabled = true
	m.stream = out
	m.bpm = DefaultBPM
	m.baseVelocity = DefaultVelocity
	m.mutex = new(sync.Mutex)
	return m, nil
}

func (m *Midi) Close() {
	if m.enabled {
		m.stream.Abort()
		m.stream.Close()
	}
	portmidi.Terminate()
}

// 34 = purple
func print(arg interface{}) {
	fmt.Printf("\033[2;34m" + fmt.Sprintf("%v ", arg) + "\033[0m")
}
