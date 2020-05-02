package midi

import (
	"math"
	"time"

	"github.com/emicklei/melrose"
)

// Play is part of melrose.AudioDevice
// It schedules all the notes on the timeline beginning at a give time (now or in the future).
// Returns the end time of the last played Note.
func (m *Midi) Play(seq melrose.Sequenceable, bpm float64, beginAt time.Time) time.Time {
	moment := beginAt
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
		if len(eachGroup) == 0 {
			continue
		}
		var actualDuration time.Duration
		var event midiEvent
		if canCombineMidiEvents(eachGroup) {
			// combined
			actualDuration = time.Duration(float32(wholeNoteDuration) * eachGroup[0].DurationFactor())
			event = m.combinedMidiEvent(channel, eachGroup)
			if m.echo {
				event.echoString = melrose.StringFromNoteGroup(eachGroup)
			}
		} else {
			// one-by-one
			for i, eachNote := range eachGroup {
				actualDuration = time.Duration(float32(wholeNoteDuration) * eachNote.DurationFactor())
				if eachNote.IsRest() {
					continue
				}
				event = m.combinedMidiEvent(channel, eachGroup[i:i+1])
				if m.echo {
					event.echoString = eachNote.String()
				}
			}
		}
		m.timeline.Schedule(event, moment)
		moment = moment.Add(actualDuration)
		// note off is not echoed
		m.timeline.Schedule(event.asNoteoff(), moment)
	}
	return moment
}

// Pre: notes not empty
func (m *Midi) combinedMidiEvent(channel int, notes []melrose.Note) midiEvent {
	v := notes[0].VelocityFactor()
	velocity := int(float32(m.baseVelocity) * v)
	if velocity > 127 {
		velocity = 127
	}
	if velocity < 1 {
		velocity = DefaultVelocity
	}
	nrs := []int64{}
	for _, each := range notes {
		nrs = append(nrs, int64(each.MIDI()))
	}
	return midiEvent{
		which:    nrs,
		onoff:    noteOn,
		channel:  channel,
		velocity: int64(velocity),
		out:      m.stream,
	}
}

func canCombineMidiEvents(notes []melrose.Note) bool {
	if len(notes) < 2 {
		return false
	}
	d := notes[0].DurationFactor()
	v := notes[0].VelocityFactor()
	for _, each := range notes[1:] {
		if each.DurationFactor() != d {
			return false
		}
		if each.VelocityFactor() != v {
			return false
		}
	}
	return true
}
