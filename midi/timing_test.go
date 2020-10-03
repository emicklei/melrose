package midi

import (
	"testing"
)

func TestTimingRandomOffset_NoteOffsets(t *testing.T) {
	tim := newTimingOffset(-1, 3, -2, 4)
	{
		off := tim.NoteOff()
		t.Log("off", off)
	}
}
