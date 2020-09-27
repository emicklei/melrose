package midi

import (
	"fmt"
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
	if m.echo {
		fmt.Println() // start new line
	}
	channel := m.defaultOutputChannel
	if sel, ok := seq.(core.ChannelSelector); ok {
		channel = sel.Channel()
	}
	actualSequence := seq.S()
	wholeNoteDuration := core.WholeNoteDuration(bpm)
	for _, eachGroup := range actualSequence.Notes {
		if len(eachGroup) == 0 {
			continue
		}
		if m.handledPedalChange(channel, m.timeline, moment, eachGroup) {
			continue
		}
		var actualDuration time.Duration
		var event midiEvent
		if canCombineMidiEvents(eachGroup) {
			// combined, first note makes duration and velocity
			actualDuration = time.Duration(float32(wholeNoteDuration) * eachGroup[0].DurationFactor())
			event = m.combinedMidiEvent(channel, eachGroup)
			event.echoString = core.StringFromNoteGroup(eachGroup)
		} else {
			// one-by-one
			for i, eachNote := range eachGroup {
				actualDuration = time.Duration(float32(wholeNoteDuration) * eachNote.DurationFactor())
				if eachNote.IsRest() {
					event.echoString = eachNote.String()
					continue
				}
				event = m.combinedMidiEvent(channel, eachGroup[i:i+1])
				event.echoString = eachNote.String()
			}
		}
		m.timeline.Schedule(event, moment)
		moment = moment.Add(actualDuration)
		m.timeline.Schedule(event.asNoteoff(), moment)
	}
	return moment
}

// Pre: notes not empty
func (m *Midi) combinedMidiEvent(channel int, notes []core.Note) midiEvent {
	// first note makes duration and velocity
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
		out:      m.stream,
	}
}

func canCombineMidiEvents(notes []core.Note) bool {
	// assumes group of notes does not have =<>^
	return len(notes) >= 2
}
