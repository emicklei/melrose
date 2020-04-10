package midi

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/emicklei/melrose"
)

type ChannelSelector struct {
	Target melrose.Sequenceable
	Number int
}

func (c ChannelSelector) S() melrose.Sequence {
	return c.Target.S()
}

func Channel(s melrose.Sequenceable, nr int) ChannelSelector {
	return ChannelSelector{Target: s, Number: nr}
}

func (m *Midi) Play(seq melrose.Sequenceable, echo bool) {
	if !m.enabled {
		return
	}
	channel := 1 // default
	if ch, ok := seq.(ChannelSelector); ok {
		channel = ch.Number
	}
	actualSequence := seq.S()
	wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/m.bpm))) * time.Millisecond
	for _, eachGroup := range actualSequence.Notes {
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
				m.playNote(channel, int(float32(m.baseVelocity)*n.VelocityFactor()), n, wholeNoteDuration)
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
	if velocity > 127 {
		velocity = 127
	}
	if velocity < 0 {
		velocity = DefaultVelocity
	}
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
