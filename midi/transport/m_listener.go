package transport

import (
	"slices"
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
	// bpm provides the tempo used to quantize the length of played notes
	bpm func() float64
}

const defaultBPM = 120.0

func newMListener() *mListener {
	return &mListener{
		listening:     false,
		noteOn:        map[int]mNoteEvent{},
		noteListeners: []core.NoteListener{},
		keyListeners:  map[int]core.NoteListener{},
		bpm:           func() float64 { return defaultBPM },
	}
}

// SetBPMProvider sets the function that returns the current tempo.
func (l *mListener) SetBPMProvider(bpm func() float64) {
	if bpm == nil {
		return
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.bpm = bpm
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
	ch := int(int16(0x0F)&status) + 1

	// controlChange before noteOn
	if (status & controlChange) == controlChange {
		l.mutex.RLock()
		listeners := slices.Clone(l.noteListeners)
		l.mutex.RUnlock()
		for _, each := range listeners {
			each.ControlChange(ch, nr, int(data2))
		}
		return
	}
	isNoteOn := (status & noteOn) == noteOn
	velocity := data2
	if isNoteOn && velocity > 0 {
		onNote, _ := core.MIDItoNote(0.25, nr, velocity) // length is computed on note off
		// noteOn is mutated, so a write lock is required
		l.mutex.Lock()
		if _, ok := l.noteOn[nr]; ok {
			l.mutex.Unlock()
			return
		}
		l.noteOn[nr] = mNoteEvent{
			note: onNote,
			when: time.Now(),
		}
		listeners := slices.Clone(l.noteListeners)
		keyHandler, hasKeyHandler := l.keyListeners[nr]
		l.mutex.Unlock()

		notify.Debugf("on  %s", onNote)
		// listeners are called without the lock because they may add or remove listeners
		for _, each := range listeners {
			each.NoteOn(ch, onNote)
		}
		// notify key listeners
		if hasKeyHandler {
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
		l.mutex.Lock()
		on, ok := l.noteOn[nr]
		if !ok {
			l.mutex.Unlock()
			return
		}
		delete(l.noteOn, nr)
		bpm := l.bpm()
		listeners := slices.Clone(l.noteListeners)
		keyHandler, hasKeyHandler := l.keyListeners[nr]
		l.mutex.Unlock()

		// compute delta
		nanos := time.Since(on.when)
		frac, dotted := core.DurationToFraction(bpm, nanos)
		name, octave, accidental := core.MIDIToNoteParts(nr)
		offNote, _ := core.NewNote(name, octave, frac, accidental, dotted, on.note.Velocity)
		notify.Debugf("off %s [%d]", offNote, nanos)
		for _, each := range listeners {
			each.NoteOff(ch, offNote)
		}
		if hasKeyHandler {
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
