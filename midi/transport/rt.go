package transport

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"gitlab.com/gomidi/rtmididrv/imported/rtmidi"
)

func init() { Initializer = rtInitialize }

func rtInitialize() {
	if core.IsDebug() {
		notify.Debugf("transport.init: use RtmidiTransporter")
	}
	Factory = func() Transporter {
		return RtmidiTransporter{}
	}
}

type RtmidiTransporter struct{}

func (t RtmidiTransporter) HasInputCapability() bool {
	return true
}
func (t RtmidiTransporter) PrintInfo(inID, outID int) {
	notify.PrintHighlighted("usage:")
	fmt.Println(":m echo                               --- toggle printing the notes that are send")
	fmt.Println(":m in      <device-id>                --- change the default MIDI input  device id")
	fmt.Println(":m out     <device-id>                --- change the default MIDI output device id")
	fmt.Println(":m channel <device-id> <midi-channel> --- change the default MIDI channel for an output device id")
	fmt.Println()

	notify.PrintHighlighted("available:")

	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		log.Fatalln("can't open default MIDI in: ", err)
	}
	defer in.Close()
	ports, err := in.PortCount()
	if err != nil {
		log.Fatalln("can't get number of in ports: ", err.Error())
	}
	for i := 0; i < ports; i++ {
		name, err := in.PortName(i)
		if err != nil {
			name = ""
		}
		fmt.Printf("[midi] input device %d: %s\n", i, name)
	}
	{
		// Outs
		out, err := rtmidi.NewMIDIOutDefault()
		if err != nil {
			log.Fatalln("can't open default MIDI out: ", err)
		}
		defer out.Close()
		ports, err := out.PortCount()
		if err != nil {
			log.Fatalln("can't get number of out ports: ", err.Error())
		}

		for i := 0; i < ports; i++ {
			name, err := out.PortName(i)
			if err != nil {
				name = ""
			}
			fmt.Printf("[midi] output device %d: %s\n", i, name)
		}
	}
	fmt.Println()
}
func (t RtmidiTransporter) DefaultOutputDeviceID() int {
	return 0
}
func (t RtmidiTransporter) DefaultInputDeviceID() int {
	return 0
}

func (t RtmidiTransporter) NewMIDIOut(id int) (MIDIOut, error) {
	out, err := rtmidi.NewMIDIOutDefault()
	if err != nil {
		return nil, err
	}
	err = out.OpenPort(id, "")
	if err != nil {
		return nil, err
	}
	return RtmidiOut{out: out}, nil
}
func (t RtmidiTransporter) NewMIDIIn(id int) (MIDIIn, error) {
	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		return nil, err
	}
	err = in.OpenPort(id, "")
	if err != nil {
		return nil, err
	}
	return RtmidiIn{in: in}, nil
}
func (t RtmidiTransporter) Terminate() {
	// noop
}
func (t RtmidiTransporter) NewMIDIListener(in MIDIIn) MIDIListener {
	return newRtListener(in.(RtmidiIn).in)
}

type RtmidiOut struct {
	out rtmidi.MIDIOut
}

func (o RtmidiOut) WriteShort(status int64, data1 int64, data2 int64) error {
	return o.out.SendMessage([]byte{byte(status & 0xFF), byte(data1 & 0xFF), byte(data2 & 0xFF)})
}
func (o RtmidiOut) Close() error { return o.out.Close() }
func (o RtmidiOut) Abort() error { return nil }

type RtmidiIn struct {
	in rtmidi.MIDIIn
}

func (i RtmidiIn) Close() error { return i.in.Close() }

type rtNoteEvent struct {
	note core.Note
	when time.Time
}
type RtListener struct {
	mutex *sync.RWMutex

	listening     bool
	midiIn        rtmidi.MIDIIn
	noteOn        map[int]rtNoteEvent
	noteListeners []core.NoteListener
	keyListeners  map[int]core.NoteListener
}

func newRtListener(in rtmidi.MIDIIn) *RtListener {
	return &RtListener{
		midiIn:       in,
		mutex:        new(sync.RWMutex),
		noteOn:       map[int]rtNoteEvent{},
		keyListeners: map[int]core.NoteListener{},
	}
}

func (l *RtListener) Add(lis core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.noteListeners = append(l.noteListeners, lis)
}

func (l *RtListener) Remove(lis core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	without := []core.NoteListener{}
	for _, each := range l.noteListeners {
		if each != lis {
			without = append(without, each)
		}
	}
	l.noteListeners = without
}
func (l *RtListener) OnKey(note core.Note, handler core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	nr := note.MIDI()
	// remove existing for the key
	old, ok := l.keyListeners[nr]
	if ok {
		l.Remove(old)
		delete(l.keyListeners, nr)
	}
	if handler == nil {
		return
	}
	// add to map and list
	l.keyListeners[nr] = handler
	l.noteListeners = append(l.noteListeners, handler)
}

func (l *RtListener) Start() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.listening {
		return
	}
	l.listening = true
	// since l.midiIn.SetCallback is blocking on success, there is no meaningful way to get an error
	// and set the callback non blocking
	go func() {
		if err := l.midiIn.SetCallback(l.handleEvent); err != nil {
			notify.Warnf("failed to set listener callback")
		}
	}()
}

// data = status,data1,data2
func (l *RtListener) handleEvent(m rtmidi.MIDIIn, data []byte, delta float64) {
	if core.IsDebug() {
		notify.Debugf("transport.RtListener.handleEvent data=%v,f=%v", data, delta)
	}
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	status := int64(data[0])
	nr := int(data[1])
	data2 := int(data[2])

	// controlChange before noteOn
	isControlChange := (status & controlChange) == controlChange
	if isControlChange {
		for _, each := range l.noteListeners {
			// TODO get channel
			each.ControlChange(0, nr, int(data2))
		}
		return
	}
	isNoteOn := (status & noteOn) == noteOn
	velocity := data2
	if isNoteOn && velocity > 0 {
		if _, ok := l.noteOn[nr]; ok {
			return
		}
		onNote, _ := core.MIDItoNote(0.25, nr, velocity)
		l.noteOn[nr] = rtNoteEvent{
			note: onNote,
			when: time.Now(),
		}
		for _, each := range l.noteListeners {
			each.NoteOn(onNote)
		}
		return
	}
	isNoteOff := (status & noteOff) == noteOff
	// for devices that support aftertouch, a noteOn with velocity 0 is also handled as a noteOff
	if !isNoteOff {
		isNoteOff = isNoteOn && velocity == 0
	}
	if isNoteOff {
		on, ok := l.noteOn[nr]
		if !ok {
			return
		}
		delete(l.noteOn, nr)
		// compute delta
		ms := time.Duration(time.Now().UnixNano()-on.when.UnixNano()) * time.Nanosecond
		frac := core.DurationToFraction(120.0, ms) // TODO
		offNote, _ := core.MIDItoNote(frac, nr, core.Normal)
		for _, each := range l.noteListeners {
			each.NoteOff(offNote)
		}
		return
	}
}
func (l *RtListener) Stop() {
	if !l.listening {
		return
	}
	if err := l.midiIn.CancelCallback(); err != nil {
		notify.Warnf("failed to cancel listener callback")
	}
}
