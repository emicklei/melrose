package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi/transport"
	"github.com/emicklei/melrose/notify"
)

// TODO deprecated
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
func (registry *DeviceRegistry) Play(condition core.Condition, seq core.Sequenceable, bpm float64, beginAt time.Time) time.Time {
	if core.IsDebug() {
		notify.Debugf("midi.play: time=%s object=%s", beginAt.Format("04:05.000"), core.Storex(seq))
	}
	// unwrap if variable because we need to detect device or channel selector
	seq = core.UnValue(seq)

	// which device?
	var device *OutputDevice
	deviceID := registry.defaultOutputID
	if dev, ok := seq.(core.DeviceSelector); ok {
		deviceID = dev.DeviceID()
		seq = dev.Unwrap()
	}
	device, err := registry.Output(deviceID)
	if err != nil {
		return beginAt
	}

	// which channel?
	channel := device.defaultChannel
	if sel, ok := seq.(core.ChannelSelector); ok {
		channel = sel.Channel()
		seq = sel.Unwrap()
	}

	// schedule all notes of the sequenceable
	wholeNoteDuration := core.WholeNoteDuration(bpm)
	moment := beginAt
	for _, eachGroup := range seq.S().Notes {
		if len(eachGroup) == 0 {
			continue
		}
		// pedal
		if device.handledPedalChange(condition, channel, device.timeline, moment, eachGroup) {
			continue
		}
		// one note
		if len(eachGroup) == 1 {
			moment = scheduleOneNote(device, condition, channel, eachGroup[0], wholeNoteDuration, moment)
			continue
		}
		//  more than one note
		if canCombineEvent(eachGroup) {
			event := combinedMidiEvent(channel, eachGroup, device.stream)
			if device.echo {
				event.echoString = core.StringFromNoteGroup(eachGroup)
			}
			actualDuration := durationOfGroup(eachGroup, wholeNoteDuration)
			event.mustHandle = condition
			moment = scheduleOnOffEvents(device, event, actualDuration, moment)
			continue
		}
		//  not combinable group of more than one note
		earliest := moment.Add(1 * time.Hour)
		for _, each := range eachGroup {
			endTime := scheduleOneNote(device, condition, channel, each, wholeNoteDuration, moment)
			if endTime.Before(earliest) {
				earliest = endTime
			}
		}
		moment = earliest
	}
	return moment
}

// returns the longest TODO in core?
func durationOfGroup(notes []core.Note, whole time.Duration) time.Duration {
	longest := time.Duration(0)
	for _, each := range notes {
		eachDuration := time.Duration(float32(whole) * each.DurationFactor())
		if eachDuration > longest {
			longest = eachDuration
		}
	}
	return longest
}

func scheduleOneNote(device *OutputDevice, condition core.Condition, channel int, note core.Note, whole time.Duration, moment time.Time) time.Time {
	if note.IsRest() {
		event := restEvent{mustHandle: condition}
		if device.echo {
			event.echoString = note.String()
		}
		device.timeline.Schedule(event, moment)
		actualDuration := time.Duration(float32(whole) * note.DurationFactor())
		return moment.Add(actualDuration)
	}
	// midi variable length note?
	if fixed, ok := note.NonFractionBasedDuration(); ok {
		event := midiEvent{
			which:      []int64{int64(note.MIDI())},
			onoff:      noteOn,
			channel:    channel,
			velocity:   int64(note.Velocity),
			out:        device.stream,
			mustHandle: condition,
		}
		return scheduleOnOffEvents(device, event, fixed, moment)
	}
	// normal note
	event := midiEvent{
		which:      []int64{int64(note.MIDI())},
		onoff:      noteOn,
		channel:    channel,
		velocity:   int64(note.Velocity),
		out:        device.stream,
		mustHandle: condition,
	}
	actualDuration := time.Duration(float32(whole) * note.DurationFactor())
	return scheduleOnOffEvents(device, event, actualDuration, moment)

}

func scheduleOnOffEvents(device *OutputDevice, event midiEvent, duration time.Duration, at time.Time) time.Time {
	device.timeline.Schedule(event, at)
	moment := at.Add(duration)
	device.timeline.Schedule(event.asNoteoff(), moment)
	return moment
}

func canCombineEvent(notes []core.Note) bool {
	if len(notes) <= 1 {
		return true
	}
	dur, vel := notes[0].DurationFactor(), notes[0].Velocity
	for n := 1; n < len(notes); n++ {
		d, v := notes[n].DurationFactor(), notes[n].Velocity
		if d != dur || v != vel {
			return false
		}
	}
	return true
}

// Pre: notes not empty
func combinedMidiEvent(channel int, notes []core.Note, stream transport.MIDIOut) midiEvent {
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
