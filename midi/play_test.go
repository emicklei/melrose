package midi

import (
	"math"
	"testing"
	"time"

	"github.com/emicklei/melrose"
)

func TestDurations(t *testing.T) {
	for _, bpm := range []float64{60, 120, 240, 300} {
		t.Log("bpm", bpm)
		wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/bpm))) * time.Millisecond
		t.Log("whole", wholeNoteDuration)
		s := melrose.S("1C 2C 4C 8C 16C")
		s.NotesDo(func(each melrose.Note) {
			actualDuration := time.Duration(float32(wholeNoteDuration) * each.DurationFactor())
			t.Log(each.String(), actualDuration)
		})
		t.Log("-----")
	}
}
