package term

type StringChangeHandler func(old, new string)
type StringListChangeHandler func(old, new []string)

type StringHolder struct {
	Value      string
	dependents []StringChangeHandler
}

func (s *StringHolder) AddDependent(h StringChangeHandler) {
	s.dependents = append(s.dependents, h)
}

func (s *StringHolder) Set(newValue string) {
	old := s.Value
	if newValue == old {
		return
	}
	s.Value = newValue
	for _, each := range s.dependents[:] {
		each(old, newValue)
	}
}

type StringListSelectionHolder struct {
	Selection      string
	SelectionIndex int
	List           []string
	dependents     []StringListChangeHandler
}

func (s *StringListSelectionHolder) Set(newValue []string) {
	old := s.List
	if len(old) == len(newValue) {
		same := true
		for i, each := range old {
			if same && newValue[i] != each {
				same = false
			}
		}
		if same {
			return
		}
	}
	s.List = newValue
	for _, each := range s.dependents[:] {
		each(old, newValue)
	}
}

func (s *StringListSelectionHolder) AddDependent(h StringListChangeHandler) {
	s.dependents = append(s.dependents, h)
}

type WriterStringHolderAdaptor struct {
	target *StringHolder
}

func (w WriterStringHolderAdaptor) Write(data []byte) (int, error) {
	w.target.Set(w.target.Value + string(data))
	return len(data), nil
}
