package core

// NoteGroup is an ordered list of Note with special playing characteristics
type NoteGroup struct {
	Notes      []Note
	isParallel bool
	isSustain  bool
}

// NewGroup returns a NoteGroup value
func NewGroup(n []Note, isParallel bool, isSustain bool) NoteGroup {
	return NoteGroup{
		Notes:      n,
		isParallel: isParallel,
		isSustain:  isSustain,
	}
}

func (g NoteGroup) First() Note {
	return g.Notes[0]
}
