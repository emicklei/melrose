package midi

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/emicklei/melrose"
	"github.com/rakyll/portmidi"
)

type Midi struct {
	enabled         bool
	stream          *portmidi.Stream
	mutex           *sync.Mutex
	deviceID        int
	echo            bool
	bpm             float64
	defaultVelocity int
}

const (
	noteOn  int = 0x90
	noteOff int = 0x80
)

func (m *Midi) Play(seq melrose.Sequence, echo bool) {
	if !m.enabled {
		return
	}
	wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/m.bpm))) * time.Millisecond
	for _, eachGroup := range seq.Notes {
		if echo {
			if len(eachGroup) == 1 {
				print(eachGroup[0])
			} else {
				print(eachGroup)
			}
		}
		wg := new(sync.WaitGroup)
		for _, eachNote := range eachGroup {
			wg.Add(1)
			go func(n melrose.Note) {
				m.playNote(1, int(float32(m.defaultVelocity)*n.VelocityFactor()), n, wholeNoteDuration)
				wg.Done()
			}(eachNote)
		}
		wg.Wait()
	}
	if echo {
		fmt.Println()
	}
}

func (m *Midi) playNote(channel int, velocity int, note melrose.Note, wholeNoteDuration time.Duration) {
	actualDuration := time.Duration(float32(wholeNoteDuration) * note.DurationFactor())
	if note.IsRest() {
		time.Sleep(actualDuration)
		return
	}
	data1 := note.MIDI()

	//fmt.Println("on", data1, actualDuration)
	m.mutex.Lock()
	m.stream.WriteShort(int64(noteOn|channel), int64(data1), int64(velocity))
	m.mutex.Unlock()

	time.Sleep(actualDuration)

	//fmt.Println("off", data1)
	m.mutex.Lock()
	m.stream.WriteShort(int64(noteOff|channel), int64(data1), 100)
	m.mutex.Unlock()
}

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
	fmt.Println("[midi] default:", defaultOut)
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
	m.bpm = 120
	m.defaultVelocity = 50
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
