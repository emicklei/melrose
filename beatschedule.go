package melrose

import "sync"

type BeatSchedule struct {
	mutex   sync.RWMutex
	entries map[int64][]BeatAction
}

type BeatAction func(beat int64)

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
	return actions
}
