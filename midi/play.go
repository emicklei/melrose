package midi

import (
	"math"
	"time"

	"github.com/emicklei/melrose"
)

// Play is part of melrose.AudioDevice
// It schedules all the notes on the timeline.
// Returns the end time of the last played Note.
func (m *Midi) Play(seq melrose.Sequenceable, bpm float64) time.Time {
	moment := time.Now()
	if !m.enabled {
		return moment
	}
	channel := m.defaultOutputChannel
	if sel, ok := seq.(melrose.ChannelSelector); ok {
		channel = sel.Channel()
	}
	actualSequence := seq.S()
	wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/bpm))) * time.Millisecond // 4 = signature
	for _, eachGroup := range actualSequence.Notes {
		var actualDuration time.Duration
		for _, eachNote := range eachGroup {
			// TODO all have the same duration so combine the event
			actualDuration = time.Duration(float32(wholeNoteDuration) * eachNote.DurationFactor())
			if eachNote.IsRest() {
				continue
			}
			velocity := int(float32(m.baseVelocity) * eachNote.VelocityFactor())
			if velocity > 127 {
				velocity = 127
			}
			if velocity < 1 {
				velocity = DefaultVelocity
			}
			event := midiEvent{
				which:    []int64{int64(eachNote.MIDI())},
				onoff:    noteOn,
				channel:  channel,
				velocity: int64(velocity),
				out:      m.stream,
			}
			if m.echo {
				// only for ON
				event.echoString = eachNote.String()
			}
			m.timeline.Schedule(event, moment)
			m.timeline.Schedule(event.asNoteoff(), moment.Add(actualDuration))
		}
		moment = moment.Add(actualDuration)
	}
	return moment
}
