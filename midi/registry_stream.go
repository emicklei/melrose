package midi

import (
	"sync"

	"github.com/rakyll/portmidi"
)

type streamRegistry struct {
	mutex *sync.RWMutex
	out   map[int]MIDIOut
	in    map[int]MIDIIn
}

type MIDIOut interface {
	WriteShort(status int64, data1 int64, data2 int64) error
	Close() error
	Abort() error
}

type MIDIIn interface {
	Close() error
}

func newStreamRegistry() *streamRegistry {
	return &streamRegistry{
		mutex: new(sync.RWMutex),
		out:   map[int]MIDIOut{},
		in:    map[int]MIDIIn{},
	}
}

func (s *streamRegistry) output(id int) (MIDIOut, error) {
	s.mutex.RLock()
	if m, ok := s.out[id]; ok {
		s.mutex.RUnlock()
		return m, nil
	}
	s.mutex.RUnlock()
	// not present
	s.mutex.Lock()
	defer s.mutex.Unlock()
	out, err := portmidi.NewOutputStream(portmidi.DeviceID(id), 1024, 0) // TODO flag
	if err != nil {
		return nil, err
	}
	// TEMP TODO
	tout := tracingMIDIStreamOn(out)
	s.out[id] = tout
	return tout, nil
}

func (s *streamRegistry) input(id int) (MIDIIn, error) {
	s.mutex.RLock()
	if m, ok := s.in[id]; ok {
		s.mutex.RUnlock()
		return m, nil
	}
	s.mutex.RUnlock()
	// not present
	s.mutex.Lock()
	defer s.mutex.Unlock()
	in, err := portmidi.NewInputStream(portmidi.DeviceID(id), 1024) // TODO flag
	if err != nil {
		return nil, err
	}
	s.in[id] = in
	return in, nil
}

func (s *streamRegistry) close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, each := range s.in {
		each.Close()
	}
	for _, each := range s.out {
		each.Abort()
		each.Close()
	}
	s.out = map[int]MIDIOut{}
	s.in = map[int]MIDIIn{}
	return nil
}
