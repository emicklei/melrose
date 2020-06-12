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
			actualDuration := time.Duration(float32(wholeNoteDuration) * each.Length())
			t.Log(each.String(), actualDuration)
		})
		t.Log("-----")
	}
}

func TestEventNoteOff(t *testing.T) {
	on := midiEvent{onoff: noteOn}
	off := on.asNoteoff()
	if got, want := on.onoff, noteOn; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := off.onoff, noteOff; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestPlay(t *testing.T) {
	m := new(Midi)
	m.enabled = true
	m.timeline = melrose.NewTimeline()
	now := time.Now()
	m.Play(melrose.MustParseSequence("C D"), 120.0, now)
	m.timeline.EventsDo(func(event melrose.TimelineEvent, when time.Time) {
		t.Logf("on [%v] event [%v]\n", when.Sub(now), event)
	})
}
