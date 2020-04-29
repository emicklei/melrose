package melrose

import (
	"testing"
	"time"
)

func TestRecorderAdd(t *testing.T) {
	r := NewRecording()
	now := time.Now()
	cOn := NewNoteChange(true, 60, 70)
	r.Add(cOn, now)
	cOff := NewNoteChange(false, 60, 70)
	r.Add(cOff, now.Add(500*time.Millisecond))
	if got, want := r.timeline.Len(), int64(2); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	s := r.S()
	t.Log(s)
}
