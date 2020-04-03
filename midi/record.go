package midi

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose"
	"github.com/rakyll/portmidi"
)

func (m *Midi) Record(deviceID int, stopAfterInactivity time.Duration) (melrose.Sequence, error) {
	in, err := portmidi.NewInputStream(portmidi.DeviceID(deviceID), 1024) // buffer
	if err != nil {
		return melrose.Sequence{}, err
	}
	defer in.Close()

	midiDeviceInfo := portmidi.Info(portmidi.DeviceID(deviceID))
	info(fmt.Sprintf("listening to %s/%s ...\n", midiDeviceInfo.Interface, midiDeviceInfo.Name))

	// listing on all channels TODO
	noteMap := map[int64]portmidi.Event{}

	notes := []melrose.Note{}
	ch := in.Listen()
	timeout := time.NewTimer(stopAfterInactivity)
	for {
		timeout.Reset(stopAfterInactivity)
		select {
		case each := <-ch:
			if each.Status == 0x90 {
				noteMap[each.Data1] = each
				continue
			}
			if each.Status != 0x80 {
				continue
			}
			startEvent := noteMap[each.Data1]
			//fmt.Println("ts,note,velocity", startEvent.Timestamp, startEvent.Data1, startEvent.Data2)
			//fmt.Println("ts,note,velocity", each.Timestamp, each.Data1, each.Data2)

			note := m.eventToNote(startEvent, each)

			//fmt.Println(startEvent.Data2, note.VelocityFactor()*float32(m.baseVelocity))

			print(note)
			notes = append(notes, note)

			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			goto done
		}
	}
done:
	info(fmt.Sprintf("\nstopped after %v of inactivity\n", stopAfterInactivity))
	return melrose.BuildSequence(notes), nil
}

// TODO compute duration
func (m *Midi) eventToNote(start, end portmidi.Event) melrose.Note {
	factor := float32(start.Data2) / float32(m.baseVelocity)
	return melrose.MIDItoNote(int(start.Data1), factor)
}
