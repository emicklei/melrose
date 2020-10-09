package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
)

// Play is part of melrose.AudioDevice
// It schedules all the notes on the timeline beginning at a give time (now or in the future).
// Returns the end time of the last played Note.
func (m *Midi) Play(seq core.Sequenceable, bpm float64, beginAt time.Time) time.Time {
	moment := beginAt
	if !m.enabled {
		return moment
	}
	channel := m.defaultOutputChannel
	if sel, ok := seq.(core.ChannelSelector); ok {
		channel = sel.Channel()
	}
	wholeNoteDuration := core.WholeNoteDuration(bpm)
	for _, eachGroup := range seq.S().Notes {
		if len(eachGroup) == 0 {
			continue
		}
		if m.handledPedalChange(channel, m.timeline, moment, eachGroup) {
			continue
		}
		var actualDuration = time.Duration(float32(wholeNoteDuration) * eachGroup[0].DurationFactor())
		var event midiEvent
		if len(eachGroup) > 1 {
			// combined, first note makes fraction and velocity
			event = m.combinedMidiEvent(channel, eachGroup)
			if m.echo {
				event.echoString = core.StringFromNoteGroup(eachGroup)
			}
		} else {
			// solo note
			// rest?
			if eachGroup[0].IsRest() {
				event := restEvent{}
				if m.echo {
					event.echoString = eachGroup[0].String()
				}
				m.timeline.Schedule(event, moment)
				moment = moment.Add(actualDuration)
				continue
			}
			// midi variable length note?
			if fixed, ok := eachGroup[0].NonFractionBasedDuration(); ok {
				actualDuration = fixed
			}
			// non-rest
			event = m.combinedMidiEvent(channel, eachGroup)
			if m.echo {
				event.echoString = eachGroup[0].String()
			}
		}
		// solo or group
		m.timeline.Schedule(event, moment)
		moment = moment.Add(actualDuration)
		m.timeline.Schedule(event.asNoteoff(), moment)
	}
	return moment
}

// Pre: notes not empty
func (m *Midi) combinedMidiEvent(channel int, notes []core.Note) midiEvent {
	// first note makes fraction and velocity
	velocity := notes[0].Velocity
	if velocity > 127 {
		velocity = 127
	}
	if velocity < 1 {
		velocity = core.Normal
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
		out:      m.outputStream,
	}
}
