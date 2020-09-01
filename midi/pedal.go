package midi

import (
	"sync"
	"time"

	"github.com/emicklei/melrose/core"
)

// SustainPedal models the state of a pedal.
type SustainPedal struct {
	mutex     sync.Mutex
	Down      bool
	OpenNotes []midiEvent
}

// NewSustainPedal returns a new.
func NewSustainPedal() *SustainPedal {
	s := &SustainPedal{
		Down:      false,
		OpenNotes: []midiEvent{},
	}
	return s
}

// Reset forgets about pending events and pedal state.
func (s *SustainPedal) Reset() {
	s.mutex.Lock()
	s.Down = false
	s.OpenNotes = []midiEvent{}
	s.mutex.Unlock()
}

// Record is storing the note event to handle it on pedal up.
func (s *SustainPedal) Record(event midiEvent) {
	s.mutex.Lock()
	s.OpenNotes = append(s.OpenNotes, event)
	s.mutex.Unlock()
}

// TakeInstruction processes a pedal instruction iff the group has a pedal change.
// Returns true if it was processed.
func (s *SustainPedal) TakeInstruction(timeline *core.Timeline, moment time.Time, group []core.Note) bool {
	if len(group) == 0 || len(group) > 1 {
		return false
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if group[0] == core.PedalUpDown {
		s.scheduleNoteOff(timeline, moment)
		return true
	}
	if group[0] == core.PedalToggle {
		if s.Down {
			s.scheduleNoteOff(timeline, moment)
		}
		s.Down = !s.Down
		return true
	}
	return false
}

// scheduleNoteOff is run in mutex protection
func (s *SustainPedal) scheduleNoteOff(timeline *core.Timeline, moment time.Time) {
	for _, event := range s.OpenNotes {
		timeline.Schedule(event.asNoteoff(), moment)
	}
	s.OpenNotes = []midiEvent{}
}
