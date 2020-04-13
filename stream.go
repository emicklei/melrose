package melrose

type NotesStream interface {
	AtEnd() bool
	Next() []Note
}

var NoStream = EmptyStream{}

type EmptyStream struct{}

func (e EmptyStream) Next() []Note { return []Note{} }
func (e EmptyStream) AtEnd() bool  { return true }

type SequenceStream struct {
	Target Sequence
	index  int
}

func (s *SequenceStream) Next() []Note {
	g := s.Target.At(s.index)
	s.index++
	return g
}

func (s *SequenceStream) AtEnd() bool {
	return s.index == s.Target.Length()
}

func (s Sequence) Stream() NotesStream {
	if s.Length() == 0 {
		return NoStream
	}
	return &SequenceStream{Target: s}
}
