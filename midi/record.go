package midi

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose"
	"github.com/rakyll/portmidi"
)

// Record is part of melrose.AudioDevice
func (m *Midi) Record(deviceID int, stopAfterInactivity time.Duration) (melrose.Sequence, error) {
	in, err := portmidi.NewInputStream(portmidi.DeviceID(deviceID), 1024) // buffer
	if err != nil {
		return melrose.Sequence{}, err
	}
	defer in.Close()

	midiDeviceInfo := portmidi.Info(portmidi.DeviceID(deviceID))
	info(fmt.Sprintf("recording from %s/%s ... [until %v silence]\n", midiDeviceInfo.Interface, midiDeviceInfo.Name, stopAfterInactivity))

	rec := melrose.NewRecorder()
	ch := in.Listen()
	timeout := time.NewTimer(stopAfterInactivity)
	needsReset := false
	now := time.Now()
	for {
		if needsReset {
			timeout.Reset(stopAfterInactivity)
			needsReset = false
		}
		select {
		case each := <-ch: // depending on the device, this may not block and other events are received
			when := now.Add(time.Duration(each.Timestamp) * time.Millisecond)
			if each.Status == noteOn {
				print(melrose.MIDItoNote(int(each.Data1), 1.0))
				rec.Add(melrose.NewNoteChange(true, each.Data1, each.Data2), when)
				needsReset = true
				continue
			}
			if each.Status != noteOff {
				continue
			}
			// note off
			needsReset = true
			rec.Add(melrose.NewNoteChange(false, each.Data1, each.Data2), when)
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			goto done
		}
	}
done:
	info(fmt.Sprintf("\nstopped after %v of inactivity\n", stopAfterInactivity))
	return rec.BuildSequence(), nil
}

// TODO compute duration
func (m *Midi) eventToNote(start, end portmidi.Event) melrose.Note {
	factor := float32(start.Data2) / float32(m.baseVelocity)
	return melrose.MIDItoNote(int(start.Data1), factor)
}
