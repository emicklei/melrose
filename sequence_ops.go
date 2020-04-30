package melrose

func (s Sequence) Pitched(semitones int) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	if semitones == 0 {
		return s
	}
	groups := [][]Note{}
	for _, group := range s.Notes {
		changed := []Note{}
		for _, each := range group {
			changed = append(changed, each.Pitched(semitones))
		}
		groups = append(groups, changed)
	}
	return Sequence{groups}
}

func (s Sequence) Reversed() Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	groups := [][]Note{}
	for c := len(s.Notes) - 1; c != -1; c-- {
		groups = append(groups, s.Notes[c])
	}
	return Sequence{groups}
}
func (s Sequence) Reverse() Reverse {
	return Reverse{Target: s}
}

func (s Sequence) Repeated(howMany int) Sequence {
	groups := [][]Note{}
	for i := 0; i < howMany; i++ {
		groups = append(groups, s.Notes...)
	}
	return Sequence{Notes: groups}
}

func (s Sequence) Octaved(howMuch int) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	groups := [][]Note{}
	for _, group := range s.Notes {
		changed := []Note{}
		for _, each := range group {
			changed = append(changed, each.Octaved(howMuch))
		}
		groups = append(groups, changed)
	}
	return Sequence{groups}
}

func (s Sequence) RotatedBy(howMany int) Sequence {
	if len(s.Notes) == 0 {
		return s
	}
	direction := Right
	if howMany < 0 {
		direction = Left
		howMany = -howMany
	}
	groups := s.Notes
	for c := 0; c < howMany; c++ {
		if direction == Left {
			first := groups[0]
			groups = append(groups[1:], first)
		} else {
			last := groups[len(groups)-1]
			groups = append([][]Note{last}, groups[:len(groups)-1]...)
		}
	}
	return Sequence{groups}
}
