package midi

import (
	"math"
	"time"

	"github.com/emicklei/melrose"
	"github.com/rakyll/portmidi"
)

type Midi struct {
	enabled bool
	stream  *portmidi.Stream
	bpm     float64
}

const (
	noteOn  = 0x90
	noteOff = 0x80
)

func (m Midi) Play(seq melrose.Sequence) {
	if !m.enabled {
		return
	}
	wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/m.bpm))) * time.Millisecond
	//fmt.Println("whole", wholeNoteDuration)
	seq.NotesDo(func(each melrose.Note) {
		data1 := each.MIDI()
		actualDuration := time.Duration(float32(wholeNoteDuration) * each.DurationFactor())
		//fmt.Println("on", data1, actualDuration)
		m.stream.WriteShort(noteOn, int64(data1), 100)
		time.Sleep(actualDuration)
		//fmt.Println("off", data1)
		m.stream.WriteShort(noteOff, int64(data1), 100)
	})
}

// BeatsPerMinute (BPM) ; beats each the length of a quarter note per minute.
func (m *Midi) BeatsPerMinute(bpm float64) {
	m.bpm = bpm
}

func Open() (*Midi, error) {
	m := new(Midi)
	m.enabled = false
	portmidi.Initialize()
	deviceID := portmidi.DefaultOutputDeviceID()
	out, err := portmidi.NewOutputStream(deviceID, 1024, 0)
	if err != nil {
		return nil, err
	}
	m.enabled = true
	m.stream = out
	m.bpm = 120
	return m, nil
}

func (m *Midi) Close() {
	if m.enabled {
		m.stream.Abort()
		m.stream.Close()
	}
	portmidi.Terminate()
}
