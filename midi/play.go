package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
)

func (registry *DeviceRegistry) Schedule(e core.TimelineEvent, beginAt time.Time) {
	// TODO check DeviceSelector
	device, err := registry.Output(registry.defaultOutputID)
	if err != nil {
		return
	}
	device.timeline.Schedule(e, beginAt)
}

// Play schedules all the notes on the timeline beginning at a give time (now or in the future).
// Returns the end time of the last played Note.
func (registry *DeviceRegistry) Play(seq core.Sequenceable, bpm float64, beginAt time.Time) time.Time {
	// which device?
	var device *OutputDevice
	deviceID := registry.defaultOutputID
	if dev, ok := seq.(core.DeviceSelector); ok {
		deviceID = dev.DeviceID()
	}
	device, err := registry.Output(deviceID)
	if err != nil {
		return beginAt
	}

	// which channel?
	channel := device.defaultChannel
	if sel, ok := seq.(core.ChannelSelector); ok {
		channel = sel.Channel()
	}

	// schedule all notes of the sequenceable
	wholeNoteDuration := core.WholeNoteDuration(bpm)
	moment := beginAt
	for _, eachGroup := range seq.S().Notes {
		if len(eachGroup) == 0 {
			continue
		}
		if device.handledPedalChange(channel, device.timeline, moment, eachGroup) {
			continue
		}
		var actualDuration = time.Duration(float32(wholeNoteDuration) * eachGroup[0].DurationFactor())
		var event midiEvent
		if len(eachGroup) > 1 {
			// combined, first note makes fraction and velocity
			event = combinedMidiEvent(channel, eachGroup, device.stream)
			if device.echo {
				event.echoString = core.StringFromNoteGroup(eachGroup)
			}
		} else {
			// solo note
			// rest?
			if eachGroup[0].IsRest() {
				event := restEvent{}
				if device.echo {
					event.echoString = eachGroup[0].String()
				}
				device.timeline.Schedule(event, moment)
				moment = moment.Add(actualDuration)
				continue
			}
			// midi variable length note?
			if fixed, ok := eachGroup[0].NonFractionBasedDuration(); ok {
				actualDuration = fixed
			}
			// non-rest
			event = combinedMidiEvent(channel, eachGroup, device.stream)
			if device.echo {
				event.echoString = eachGroup[0].String()
			}
		}
		// solo or group
		device.timeline.Schedule(event, moment)
		moment = moment.Add(actualDuration)
		device.timeline.Schedule(event.asNoteoff(), moment)
	}
	return moment
}

// Pre: notes not empty
func combinedMidiEvent(channel int, notes []core.Note, stream MIDIOut) midiEvent {
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
		out:      stream,
	}
}
