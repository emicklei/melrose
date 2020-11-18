package midi

import (
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/rakyll/portmidi"
)

type listener struct {
	listening bool
	stream    *portmidi.Stream
	quit      chan bool
	noteOn    map[int]portmidi.Event

	mutex         *sync.RWMutex
	noteListeners []core.NoteListener
}

func newListener(inputStream *portmidi.Stream) *listener {
	return &listener{
		stream: inputStream,
		noteOn: map[int]portmidi.Event{},
		mutex:  new(sync.RWMutex),
	}
}

func (l *listener) add(lis core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.noteListeners = append(l.noteListeners, lis)
}

func (l *listener) remove(lis core.NoteListener) {
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

func (l *listener) start() {
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
	nr := int(event.Data1)
	if event.Status == noteOn {
		if _, ok := l.noteOn[nr]; ok {
			return
		}
		// replace with now in nanos
		event.Timestamp = portmidi.Timestamp(time.Now().UnixNano())
		l.noteOn[nr] = event
		noteOn, _ := core.MIDItoNote(0.25, nr, core.Normal) // TODO
		for _, each := range l.noteListeners {
			each.NoteOn(noteOn)
		}
	} else if event.Status == noteOff {
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
		noteOff, _ := core.MIDItoNote(frac, nr, int(on.Data2))
		for _, each := range l.noteListeners {
			each.NoteOff(noteOff)
		}
	}
}

func (l *listener) stop() {
	// forget open notes
	l.noteOn = map[int]portmidi.Event{}
	if l.listening {
		l.quit <- true
	}
}
