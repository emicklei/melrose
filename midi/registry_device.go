package midi

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

var _ core.AudioDevice = (*DeviceRegistry)(nil)

type DeviceRegistry struct {
	mutex           *sync.RWMutex
	in              map[int]*InputDevice
	out             map[int]*OutputDevice
	defaultInputID  int
	defaultOutputID int
	streamRegistry  *streamRegistry
}

func NewDeviceRegistry() (*DeviceRegistry, error) {
	r := &DeviceRegistry{
		mutex:           new(sync.RWMutex),
		in:              map[int]*InputDevice{},
		out:             map[int]*OutputDevice{},
		streamRegistry:  newStreamRegistry(),
		defaultInputID:  -1,
		defaultOutputID: -1,
	}
	if err := r.init(); err != nil {
		return nil, err
	}
	return r, nil
}

// DefaultDeviceIDs is part of AudioDevice
func (r *DeviceRegistry) DefaultDeviceIDs() (inputDeviceID, outputDeviceID int) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.defaultInputID, r.defaultOutputID
}

func (r *DeviceRegistry) Reset() {
	for _, each := range r.out {
		each.Reset()
	}
	for _, each := range r.in {
		each.stopListener()
	}
}

func (r *DeviceRegistry) Output(id int) (*OutputDevice, error) {
	if id == -1 {
		return nil, errors.New("no output available")
	}
	r.mutex.RLock()
	if m, ok := r.out[id]; ok {
		r.mutex.RUnlock()
		return m, nil
	}
	r.mutex.RUnlock()
	// not present
	r.mutex.Lock()
	defer r.mutex.Unlock()
	midiOut, err := r.streamRegistry.output(id)
	if err != nil {
		return nil, err
	}
	od := NewOutputDevice(id, midiOut, 1, core.NewTimeline())
	r.out[id] = od
	od.Start() // play outgoing notes
	return od, nil
}

func (r *DeviceRegistry) Input(id int) (*InputDevice, error) {
	if id == -1 {
		return nil, errors.New("no input available")
	}
	r.mutex.RLock()
	if m, ok := r.in[id]; ok {
		r.mutex.RUnlock()
		return m, nil
	}
	r.mutex.RUnlock()
	// not present
	r.mutex.Lock()
	defer r.mutex.Unlock()
	midiIn, err := r.streamRegistry.input(id)
	if err != nil {
		return nil, err
	}
	ide := NewInputDevice(id, midiIn, r.streamRegistry.transport)
	r.in[id] = ide
	// do not start listening until requested for
	return ide, nil
}

func (r *DeviceRegistry) init() error {
	r.defaultOutputID = r.streamRegistry.transport.DefaultOutputDeviceID()
	r.defaultInputID = r.streamRegistry.transport.DefaultInputDeviceID()
	if err := r.initInputs(); err != nil {
		return err
	}
	if err := r.initOutputs(); err != nil {
		return err
	}
	return nil
}

func (r *DeviceRegistry) Report() {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	notify.PrintHighlighted("MIDI available to melrōse:")
	for k := range slices.Sorted(maps.Keys(r.in)) {
		v := r.in[k]
		fmt.Printf(" input device %d (:m i %d) = %s\n", k, k, v.name)
	}
	fmt.Println()
	for k := range slices.Sorted(maps.Keys(r.out)) {
		v := r.out[k]
		fmt.Printf("output device %d (:m o %d) = %s\n", k, k, v.name)
	}
}

func (r *DeviceRegistry) initOutputs() error {
	out, err := rtmidi.NewMIDIOutDefault()
	if err != nil {
		return fmt.Errorf("can't open default MIDI out: %w", err)
	}
	defer out.Close()
	ports, err := out.PortCount()
	if err != nil {
		return fmt.Errorf("can't get number of output ports: %w", err)
	}
	for each := range ports {
		device, err := r.Output(each)
		if err != nil {
			continue
		}
		name, err := out.PortName(each)
		if err != nil {
			name = ""
		}
		device.name = name
	}
	return nil
}

func (r *DeviceRegistry) initInputs() error {
	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		return fmt.Errorf("can't open default MIDI in: %w", err)
	}
	defer in.Close()
	ports, err := in.PortCount()
	if err != nil {
		return fmt.Errorf("can't get number of input ports: %w", err)
	}
	for each := range ports {
		device, err := r.Input(each)
		if err != nil {
			continue
		}
		name, err := in.PortName(each)
		if err != nil {
			name = ""
		}
		device.name = name
	}
	return nil
}

func (r *DeviceRegistry) Close() error {
	for _, each := range r.in {
		each.stopListener()
	}
	return r.streamRegistry.close()
}

func (r *DeviceRegistry) HasInputCapability() bool {
	return r.streamRegistry.transport.HasInputCapability()
}

func (r *DeviceRegistry) OnKey(ctx core.Context, deviceID int, channel int, note core.Note, fun core.HasValue) error {
	in, err := r.Input(deviceID)
	if err != nil {
		return fmt.Errorf("input creation failed:%v", err)
	}
	in.listener.Start()
	if fun == nil {
		in.listener.OnKey(note, nil)
		return nil
	}
	trigger := NewKeyTrigger(ctx, channel, note, fun)
	in.listener.OnKey(note, trigger)
	return nil
}

func (r *DeviceRegistry) Listen(deviceID int, who core.NoteListener, isStart bool) {
	notify.Debugf("midi.listen id=%d, start=%v", deviceID, isStart)

	in, err := r.Input(deviceID)
	if err != nil {
		notify.Warnf("input creation failed:%v", err)
		return
	}
	if isStart {
		in.listener.Start()
		// wait for pending events to be flushed
		time.Sleep(200 * time.Millisecond)
		in.listener.Add(who)
	} else {
		in.listener.Remove(who)
		// do not stop the listener ; incoming events are just ignored. otherwise buffer will overflow
	}
}

func (r *DeviceRegistry) ListDevices() (list []core.DeviceDescriptor) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for k, v := range r.in {
		list = append(list, core.DeviceDescriptor{
			ID:      k,
			IsInput: true,
			Name:    v.name,
		})
	}
	for k, v := range r.out {
		list = append(list, core.DeviceDescriptor{
			ID:      k,
			IsInput: false,
			Name:    v.name,
		})
	}
	return list

}
