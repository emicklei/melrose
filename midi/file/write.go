package file

import (
	"bufio"
	"encoding/binary"
	"log"
	"math"
	"os"
	"time"

	"github.com/Try431/EasyMIDI/smf"
	"github.com/Try431/EasyMIDI/smfio"
	"github.com/emicklei/melrose"
)

const ticksPerBeat uint16 = 960

// Export creates (overwrites) a SMF Midi file
func Export(fileName string, seq melrose.Sequenceable, bpm float64) error {
	// Create division
	// https://www.recordingblogs.com/wiki/time-division-of-a-midi-file
	division, err := smf.NewDivision(ticksPerBeat, smf.NOSMTPE)
	if err != nil {
		return err
	}

	// Create new midi struct
	midi, err := smf.NewSMF(smf.Format0, *division)
	if err != nil {
		return err
	}

	// Create track struct
	track := &smf.Track{}

	// Add track to new midi struct
	err = midi.AddTrack(track)
	if err != nil {
		return err
	}

	// https://www.recordingblogs.com/wiki/midi-set-tempo-meta-message
	// time = 10000 * (500ms / 960)  5.2 sec

	quarterMS := quarterUSFromBPM(bpm)
	tempoData := make([]byte, 4)
	binary.BigEndian.PutUint32(tempoData, quarterMS)
	tempo, err := smf.NewMetaEvent(0, smf.MetaSetTempo, tempoData[1:]) // take 3 bytes only
	if err != nil {
		return err
	}
	err = track.AddEvent(tempo)
	if err != nil {
		return err
	}

	// All the notes
	wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/bpm))) * time.Millisecond // 4 = signature TODO create func
	var moment time.Duration
	var lastTicks uint32 = 0
	for _, group := range seq.S().Notes {
		if len(group) == 0 {
			continue
		}
		channel := uint8(0x00)
		actualDuration := time.Duration(float32(wholeNoteDuration) * group[0].Length())
		if group[0].IsRest() {
			moment = moment + actualDuration
			continue
		}
		absoluteTicks := ticksFromDuration(moment, quarterMS)
		log.Println("on", moment)
		for i, each := range group {
			var deltaTicks uint32 = 0
			if i == 0 {
				deltaTicks = absoluteTicks - lastTicks
			}
			noteOn, err := smf.NewMIDIEvent(deltaTicks, smf.NoteOnStatus, channel, uint8(each.MIDI()), uint8(each.Velocity))
			if err != nil {
				return err
			}
			err = track.AddEvent(noteOn)
			if err != nil {
				return err
			}
		}
		lastTicks = absoluteTicks
		moment = moment + actualDuration
		log.Println("off", moment)
		absoluteTicks = ticksFromDuration(moment, quarterMS)
		for i, each := range group {
			var deltaTicks uint32 = 0
			if i == 0 {
				deltaTicks = absoluteTicks - lastTicks
			}
			noteOff, err := smf.NewMIDIEvent(deltaTicks, smf.NoteOffStatus, channel, uint8(each.MIDI()), 0x00) // zero velocity
			if err != nil {
				return err
			}
			err = track.AddEvent(noteOff)
			if err != nil {
				return err
			}
		}
		lastTicks = absoluteTicks
	}

	// Track end
	endTrack, err := smf.NewMetaEvent(0, smf.MetaEndOfTrack, []byte{})
	if err != nil {
		return err
	}
	err = track.AddEvent(endTrack)
	if err != nil {
		return err
	}

	// Actual write
	// Save to new midi source file
	outputMidi, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outputMidi.Close()

	// Create buffering stream
	writer := bufio.NewWriter(outputMidi)
	if err := smfio.Write(writer, midi); err != nil {
		return err
	}
	return writer.Flush()
}

func ticksFromDuration(dur time.Duration, quarterUSFromBPM uint32) uint32 {
	us := dur.Microseconds()
	return uint32(us) / quarterUSFromBPM * uint32(ticksPerBeat)
}

// duration in microseconds of one quarter note
func quarterUSFromBPM(bpm float64) uint32 {
	// 120 bpm -> 500000 usec/quarter note
	return uint32(60000000.0 / bpm)
}
