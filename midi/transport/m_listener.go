package transport

import (
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
)

type mNoteEvent struct {
	note core.Note
	when time.Time
}

type mListener struct {
	mutex sync.RWMutex

	listening     bool
	noteOn        map[int]mNoteEvent
	noteListeners []core.NoteListener
	keyListeners  map[int]core.NoteListener
}

func newMListener() *mListener {
	return &mListener{
		listening:     false,
		noteOn:        map[int]mNoteEvent{},
		noteListeners: []core.NoteListener{},
		keyListeners:  map[int]core.NoteListener{},
	}
}

func (l *mListener) Add(lis core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.noteListeners = append(l.noteListeners, lis)
}

func (l *mListener) Remove(lis core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.safeRemove(lis)
}

// safeRemove requires acquired lock
func (l *mListener) safeRemove(lis core.NoteListener) {
	without := []core.NoteListener{}
	for _, each := range l.noteListeners {
		if each != lis {
			without = append(without, each)
		}
	}
	l.noteListeners = without
}

func (l *mListener) OnKey(note core.Note, handler core.NoteListener) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	nr := note.MIDI()
	delete(l.keyListeners, nr)
	l.keyListeners[nr] = handler
}

func (l *mListener) Reset() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.noteOn = map[int]mNoteEvent{}
	l.keyListeners = map[int]core.NoteListener{}
	l.noteListeners = []core.NoteListener{}
}

func (l *mListener) HandleMIDIMessage(status int16, nr int, data2 int) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	ch := int(int16(0x0F)&status) + 1

	// controlChange before noteOn
	if (status & controlChange) == controlChange {
		for _, each := range l.noteListeners {
			each.ControlChange(ch, nr, int(data2))
		}
		return
	}
	isNoteOn := (status & noteOn) == noteOn
	velocity := data2
	if isNoteOn && velocity > 0 {
		if _, ok := l.noteOn[nr]; ok {
			return
		}
		onNote, _ := core.MIDItoNote(0.25, nr, velocity) // length is computed on note off
		l.noteOn[nr] = mNoteEvent{
			note: onNote,
			when: time.Now(),
		}
		notify.Debugf("on  %s", onNote)
		for _, each := range l.noteListeners {
			each.NoteOn(ch, onNote)
		}
		// notify key listeners
		if keyHandler, ok := l.keyListeners[nr]; ok {
			keyHandler.NoteOn(ch, onNote)
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
		frac := core.DurationToFraction(120.0, ms) // TODO, BPM
		offNote, _ := core.MIDItoNote(frac, nr, on.note.Velocity)
		notify.Debugf("off %s [%d]", offNote, ms)
		for _, each := range l.noteListeners {
			each.NoteOff(ch, offNote)
		}
		if keyHandler, ok := l.keyListeners[nr]; ok {
			keyHandler.NoteOff(ch, offNote)
		}
		return
	}
}

func (l *mListener) PrintInfo() {
	for _, each := range l.noteListeners {
		each.PrintInfo()
	}
	for _, each := range l.keyListeners {
		each.PrintInfo()
	}
}
