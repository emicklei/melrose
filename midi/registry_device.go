package midi

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/tre"
)

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
		return nil, tre.New(err, "Output", "id", id)
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
		return nil, tre.New(err, "Input", "id", id)
	}
	ide := NewInputDevice(id, midiIn, r.streamRegistry.transport)
	r.in[id] = ide
	return ide, nil
}

func (r *DeviceRegistry) init() error {
	r.defaultOutputID = r.streamRegistry.transport.DefaultOutputDeviceID()
	r.defaultInputID = r.streamRegistry.transport.DefaultInputDeviceID()
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

func (r *DeviceRegistry) Listen(deviceID int, who core.NoteListener, startOrStop bool) {
	if core.IsDebug() {
		notify.Debugf("midi.listen id=%d, start=%v", deviceID, startOrStop)
	}
	in, err := r.Input(deviceID)
	if err != nil {
		notify.Warnf("input creation failed:%v", err)
		return
	}
	if startOrStop {
		in.listener.Start()
		// wait for pending events to be flushed
		time.Sleep(200 * time.Millisecond)
		in.listener.Add(who)
	} else {
		in.listener.Remove(who)
		// do not stop the listener such that incoming events are just ignored
	}
}
