package midi

import (
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"

	"github.com/rakyll/portmidi"
)

// Record is part of melrose.AudioDevice
func (m *Device) Record(ctx core.Context) (*core.Recording, error) {
	stopAfterInactivity := time.Duration(5) * time.Second // TODO config 5
	deviceID := m.currentInputDeviceID
	return m.record(ctx, deviceID, stopAfterInactivity)
}

func (m *Device) record(ctx core.Context, deviceID int, stopAfterInactivity time.Duration) (*core.Recording, error) {
	rec := core.NewRecording()
	in, err := portmidi.NewInputStream(portmidi.DeviceID(deviceID), 1024) // buffer
	if err != nil {
		return rec, err
	}
	defer in.Close()

	midiDeviceInfo := portmidi.Info(portmidi.DeviceID(deviceID))
	fmt.Fprintf(notify.Console.StandardOut, "recording from %s/%s ... [until %v silence]\n", midiDeviceInfo.Interface, midiDeviceInfo.Name, stopAfterInactivity)

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
				rec.Add(core.NewNoteChange(true, each.Data1, each.Data2), when)
				needsReset = true
				continue
			}
			if each.Status != noteOff {
				continue
			}
			// note off
			needsReset = true
			rec.Add(core.NewNoteChange(false, each.Data1, each.Data2), when)
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			goto done
		}
	}
done:
	fmt.Fprintf(notify.Console.StandardOut, "\nstopped after %v of silence\n", stopAfterInactivity)
	core.PrintValue(ctx, rec)
	return rec, nil
}

// TODO compute duration
func (m *Device) eventToNote(start, end portmidi.Event) core.Note {
	return core.MIDItoNote(0.25, int(start.Data1), int(start.Data2)) // TODO
}
