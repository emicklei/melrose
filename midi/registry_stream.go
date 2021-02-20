package midi

import (
	"sync"

	"github.com/emicklei/melrose/midi/transport"
)

type streamRegistry struct {
	mutex     *sync.RWMutex
	out       map[int]transport.MIDIOut
	in        map[int]transport.MIDIIn
	transport transport.Transporter
}

func newStreamRegistry() *streamRegistry {
	return &streamRegistry{
		mutex:     new(sync.RWMutex),
		out:       map[int]transport.MIDIOut{},
		in:        map[int]transport.MIDIIn{},
		transport: transport.Factory(),
	}
}

func (s *streamRegistry) output(id int) (transport.MIDIOut, error) {
	s.mutex.RLock()
	if m, ok := s.out[id]; ok {
		s.mutex.RUnlock()
		return m, nil
	}
	s.mutex.RUnlock()
	// not present
	s.mutex.Lock()
	defer s.mutex.Unlock()
	out, err := s.transport.NewMIDIOut(id)
	if err != nil {
		return nil, err
	}
	s.out[id] = out
	return out, nil
}

func (s *streamRegistry) input(id int) (transport.MIDIIn, error) {
	s.mutex.RLock()
	if m, ok := s.in[id]; ok {
		s.mutex.RUnlock()
		return m, nil
	}
	s.mutex.RUnlock()
	// not present
	s.mutex.Lock()
	defer s.mutex.Unlock()
	in, err := s.transport.NewMIDIIn(id)
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
		each.Close()
	}
	s.out = map[int]transport.MIDIOut{}
	s.in = map[int]transport.MIDIIn{}
	return nil
}
