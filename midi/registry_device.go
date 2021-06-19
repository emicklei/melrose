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
func (d *DeviceRegistry) DefaultDeviceIDs() (inputDeviceID, outputDeviceID int) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.defaultInputID, d.defaultOutputID
}

func (d *DeviceRegistry) Reset() {
	for _, each := range d.out {
		each.Reset()
	}
	for _, each := range d.in {
		each.stopListener()
	}
}

func (d *DeviceRegistry) Output(id int) (*OutputDevice, error) {
	if id == -1 {
		return nil, errors.New("no output available")
	}
	d.mutex.RLock()
	if m, ok := d.out[id]; ok {
		d.mutex.RUnlock()
		return m, nil
	}
	d.mutex.RUnlock()
	// not present
	d.mutex.Lock()
	defer d.mutex.Unlock()
	midiOut, err := d.streamRegistry.output(id)
	if err != nil {
		return nil, tre.New(err, "Output", "id", id)
	}
	od := NewOutputDevice(id, midiOut, 1, core.NewTimeline())
	d.out[id] = od
	od.Start() // play outgoing notes
	return od, nil
}

func (d *DeviceRegistry) Input(id int) (*InputDevice, error) {
	if id == -1 {
		return nil, errors.New("no input available")
	}
	d.mutex.RLock()
	if m, ok := d.in[id]; ok {
		d.mutex.RUnlock()
		return m, nil
	}
	d.mutex.RUnlock()
	// not present
	d.mutex.Lock()
	defer d.mutex.Unlock()
	midiIn, err := d.streamRegistry.input(id)
	if err != nil {
		return nil, tre.New(err, "Input", "id", id)
	}
	ide := NewInputDevice(id, midiIn, d.streamRegistry.transport)
	d.in[id] = ide
	return ide, nil
}

func (d *DeviceRegistry) init() error {
	d.defaultOutputID = d.streamRegistry.transport.DefaultOutputDeviceID()
	d.defaultInputID = d.streamRegistry.transport.DefaultInputDeviceID()
	return nil
}

func (d *DeviceRegistry) Close() error {
	for _, each := range d.in {
		each.stopListener()
	}
	return d.streamRegistry.close()
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
