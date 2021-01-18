// +build !udp

package transport

import (
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/rakyll/portmidi"
)

type listener struct {
	listening bool
	stream    *portmidi.Stream
	quit      chan bool
	noteOn    map[int]portmidi.Event

	mutex         *sync.RWMutex
	noteListeners []core.NoteListener
	keyListeners  map[int]core.NoteListener
}

func newListener(inputStream *portmidi.Stream) *listener {
	return &listener{
		stream:       inputStream,
		noteOn:       map[int]portmidi.Event{},
		keyListeners: map[int]core.NoteListener{},
		mutex:        new(sync.RWMutex),
	}
}

func (l *listener) OnKey(note core.Note, handler core.NoteListener) {
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
	l.Add(handler)
}

func (l *listener) Add(lis core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.noteListeners = append(l.noteListeners, lis)
}

func (l *listener) Remove(lis core.NoteListener) {
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

func (l *listener) Start() {
	if l.listening {
		return
	}
	l.listening = true
	go l.listen()
}

func (l *listener) listen() {
	l.quit = make(chan bool)
	ch := l.stream.Listen()
	for {
		select {
		case <-l.quit:
			goto stop
		case e := <-ch:
			l.handle(e)
		}
	}
stop:
	close(l.quit)
	l.listening = false
}

func (l *listener) handle(event portmidi.Event) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if core.IsDebug() {
		notify.Debugf("portmidi.handle %v", event)
	}
	nr := int(event.Data1)
	isNoteOn := (event.Status & noteOn) == noteOn
	velocity := int(event.Data2)
	if isNoteOn && velocity > 0 {
		if _, ok := l.noteOn[nr]; ok {
			return
		}
		// replace with now in nanos
		event.Timestamp = portmidi.Timestamp(time.Now().UnixNano())
		l.noteOn[nr] = event
		onNote, _ := core.MIDItoNote(0.25, nr, velocity)
		for _, each := range l.noteListeners {
			each.NoteOn(onNote)
		}
	}
	isNoteOff := (event.Status & noteOff) == noteOff
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
		// replace with now in nanos
		event.Timestamp = portmidi.Timestamp(time.Now().UnixNano())
		// compute delta
		ms := time.Duration(event.Timestamp-on.Timestamp) * time.Nanosecond
		frac := core.DurationToFraction(120.0, ms) // TODO
		offNote, _ := core.MIDItoNote(frac, nr, velocity)
		for _, each := range l.noteListeners {
			each.NoteOff(offNote)
		}
	}
}

func (l *listener) Stop() {
	// forget open notes
	l.noteOn = map[int]portmidi.Event{}
	if l.listening {
		l.quit <- true
	}
}
