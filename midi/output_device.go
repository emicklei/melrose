package midi

import (
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/midi/transport"
	"github.com/emicklei/melrose/notify"
)

type OutputDevice struct {
	id             int
	name           string
	stream         transport.MIDIOut
	defaultChannel int

	echo     bool
	timeline *core.Timeline
}

func NewOutputDevice(id int, out transport.MIDIOut, ch int, line *core.Timeline) *OutputDevice {
	return &OutputDevice{
		id:             id,
		stream:         out,
		defaultChannel: ch,
		echo:           false,
		timeline:       line,
	}
}

func (d *OutputDevice) Start() {
	go d.timeline.Play()
}

func (d *OutputDevice) Reset() {
	defer func() {
		if err := recover(); err != nil {
			notify.Warnf("reset failed for device:%v", d.id)
		}
	}()
	// Stop the timeline goroutine first, then clear the queue
	d.timeline.Stop()
	d.timeline.Reset()
	// Give timeline time to stop processing any pending events
	time.Sleep(10 * time.Millisecond)

	if notify.IsDebug() {
		notify.Debugf("device.%d: sending Note OFF to all 16 channels", d.id)
	}
	if d.stream != nil {
		// Send a robust panic sequence for each channel.
		// Some devices respond better to CC123, others to CC120, and sustained notes
		// require a sustain-off first.
		for c := 1; c <= 16; c++ {
			status := controlChange | int64(c-1)
			notify.Debugf("reset: sending off events to [channel %d]", c)
			if err := d.stream.WriteShort(status, sustainPedal, sustainOff); err != nil {
				notify.Console.Errorf("device.%d: midi write error:%v", d.id, err)
			}
			if err := d.stream.WriteShort(status, allNotesOff, 0); err != nil {
				notify.Console.Errorf("device.%d: midi write error:%v", d.id, err)
			}
			if err := d.stream.WriteShort(status, allSoundOff, 0); err != nil {
				notify.Console.Errorf("device.%d: midi write error:%v", d.id, err)
			}
			// Small delay between channels to allow receiver to process messages.
			time.Sleep(5 * time.Millisecond)
		}

		// Send explicit note-off messages for all notes on all channels.
		// This is more direct than relying on CC123/CC120 alone and handles
		// devices/DAWs that may not respond to control changes during shutdown.
		noteOff := int64(0x80) // MIDI Note Off status
		for c := 1; c <= 16; c++ {
			status := noteOff | int64(c-1)
			notify.Debugf("reset: sending explicit note-off for all notes on channel %d", c)
			for n := 0; n <= 127; n++ {
				if err := d.stream.WriteShort(status, int64(n), 0); err != nil {
					notify.Console.Errorf("device.%d: note-off write error:%v", d.id, err)
				}
			}
			// Small delay between channels
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (d *OutputDevice) handledPedalChange(condition core.Condition, channel int, timeline *core.Timeline, moment time.Time, group []core.Note) bool {
	if len(group) == 0 || len(group) > 1 {
		return false
	}
	note := group[0]
	switch {
	case note.IsPedalUp():
		timeline.Schedule(pedalEvent{
			goingDown:  false,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
		return true
	case note.IsPedalUpDown():
		timeline.Schedule(pedalEvent{
			goingDown:  false,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
		timeline.Schedule(pedalEvent{
			goingDown:  true,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
		return true
	case note.IsPedalDown():
		timeline.Schedule(pedalEvent{
			goingDown:  true,
			channel:    channel,
			out:        d.stream,
			mustHandle: condition}, moment)
		return true
	}
	return false
}

func (d *OutputDevice) Play(condition core.Condition, seq core.Sequenceable, bpm float64, beginAt time.Time) time.Time {
	// which channel?
	channel := d.defaultChannel
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
		if d.handledPedalChange(condition, channel, d.timeline, moment, eachGroup) {
			continue
		}
		// one note
		if len(eachGroup) == 1 {
			moment = scheduleOneNote(d, condition, channel, eachGroup[0], wholeNoteDuration, moment)
			continue
		}
		//  more than one note
		if canCombineEvent(eachGroup) {
			onEvent, offEvent := combinedMidiEvents(d.id, channel, eachGroup, d.stream)
			if d.echo {
				echoStr := core.StringFromNoteGroup(eachGroup)
				onEvent.echoString = echoStr
				offEvent.echoString = echoStr
			}
			actualDuration := durationOfGroup(eachGroup, wholeNoteDuration)
			onEvent.mustHandle = condition
			offEvent.mustHandle = condition
			moment = scheduleOnOffEvents(d, onEvent, offEvent, actualDuration, moment)
			continue
		}
		//  not combinable group of more than one note
		earliest := moment.Add(1 * time.Hour)
		for _, each := range eachGroup {
			endTime := scheduleOneNote(d, condition, channel, each, wholeNoteDuration, moment)
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

	onEvent := getMidiEvent()
	onEvent.onoff = noteOn
	onEvent.device = device.id
	onEvent.channel = channel
	onEvent.velocity = int64(note.Velocity)
	onEvent.out = device.stream
	onEvent.mustHandle = condition
	onEvent.which = append(onEvent.which, int64(note.MIDI()))
	if device.echo {
		onEvent.echoString = note.String()
	}

	offEvent := getMidiEvent()
	offEvent.onoff = noteOff
	offEvent.device = device.id
	offEvent.channel = channel
	offEvent.velocity = int64(note.Velocity)
	offEvent.out = device.stream
	offEvent.mustHandle = condition
	offEvent.which = append(offEvent.which, int64(note.MIDI()))
	if device.echo {
		offEvent.echoString = note.String()
	}

	// midi variable length note?
	if fixed, ok := note.NonFractionBasedDuration(); ok {
		return scheduleOnOffEvents(device, onEvent, offEvent, fixed, moment)
	}

	actualDuration := time.Duration(float32(whole) * note.DurationFactor())
	return scheduleOnOffEvents(device, onEvent, offEvent, actualDuration, moment)
}

func scheduleOnOffEvents(device *OutputDevice, onEvent *midiEvent, offEvent *midiEvent, duration time.Duration, at time.Time) time.Time {
	device.timeline.Schedule(onEvent, at)
	moment := at.Add(duration)
	device.timeline.Schedule(offEvent, moment)
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
func combinedMidiEvents(deviceID int, channel int, notes []core.Note, stream transport.MIDIOut) (*midiEvent, *midiEvent) {
	// first note makes fraction and velocity
	velocity := notes[0].Velocity
	if velocity > 127 {
		velocity = 127
	}
	if velocity < 1 {
		velocity = core.Normal
	}

	onEvent := getMidiEvent()
	onEvent.onoff = noteOn
	onEvent.device = deviceID
	onEvent.channel = channel
	onEvent.velocity = int64(velocity)
	onEvent.out = stream

	offEvent := getMidiEvent()
	offEvent.onoff = noteOff
	offEvent.device = deviceID
	offEvent.channel = channel
	offEvent.velocity = int64(velocity)
	offEvent.out = stream

	for _, each := range notes {
		n := int64(each.MIDI())
		onEvent.which = append(onEvent.which, n)
		offEvent.which = append(offEvent.which, n)
	}

	return onEvent, offEvent
}
