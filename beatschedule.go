package melrose

import "sync"

// BeatSchedule holds mapping between beat counts and an action (function).
type BeatSchedule struct {
	mutex   *sync.RWMutex
	entries map[int64][]BeatAction
}

type BeatAction func(beat int64)

func NewBeatSchedule() *BeatSchedule {
	return &BeatSchedule{
		mutex:   new(sync.RWMutex),
		entries: map[int64][]BeatAction{},
	}
}

func (s *BeatSchedule) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.entries = map[int64][]BeatAction{}
}

func (s *BeatSchedule) Schedule(beat int64, action BeatAction) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	actions, ok := s.entries[beat]
	if ok {
		actions = append(actions, action)
	} else {
		actions = []BeatAction{action}
	}
	s.entries[beat] = actions
}

var noActions = []BeatAction{}

func (s *BeatSchedule) Unschedule(beat int64) []BeatAction {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	actions, ok := s.entries[beat]
	if !ok {
		return noActions
	}
	delete(s.entries, beat)
	return actions
}
