package midi

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/emicklei/melrose/core"
)

func TestDurations(t *testing.T) {
	for _, bpm := range []float64{60, 120, 240, 300} {
		t.Log("bpm", bpm)
		wholeNoteDuration := time.Duration(int(math.Round(4*60*1000/bpm))) * time.Millisecond
		t.Log("whole", wholeNoteDuration)
		s := core.S("1C 2C 4C 8C 16C")
		s.NotesDo(func(each core.Note) {
			actualDuration := time.Duration(float32(wholeNoteDuration) * each.DurationFactor())
			t.Log(each.String(), actualDuration)
		})
		t.Log("-----")
	}
}

func TestPlaybackTiming(t *testing.T) {
	const repetitions = 1024
	for _, bpm := range []float64{60, 120, 123, 127, 137} {
		for _, pattern := range []struct {
			name       string
			input      string
			wholeNotes int
		}{
			{"whole", "1C", 1},
			{"quarters", "C C C C", 1},
			{"sixteenths", "16C 16C 16C 16C 16C 16C 16C 16C 16C 16C 16C 16C 16C 16C 16C 16C", 1},
			{"dotted", "2.C 2.C 2.C 2.C", 3},
			{"rests", "2.= 2.= 2.= 2.=", 3},
			{"chords", "(2.C 2.E) (2.C 2.E) (2.C 2.E) (2.C 2.E)", 3},
			{"thirty-seconds", "32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C 32C", 1},
		} {
			t.Run(fmt.Sprintf("%g/%s", bpm, pattern.name), func(t *testing.T) {
				sequence := core.S(pattern.input)
				large := core.Sequence{}
				for repetition := 0; repetition < repetitions; repetition++ {
					large = large.SequenceJoin(sequence)
				}
				timeline := core.NewTimeline()
				device := NewOutputDevice(1, nil, 1, timeline)
				begin := time.Now().Add(time.Hour)
				end := core.PlayAt(device, core.NoCondition, large, bpm, begin)
				actual := end.Sub(begin)
				wholeTimeline := core.NewTimeline()
				wholeDevice := NewOutputDevice(1, nil, 1, wholeTimeline)
				wholeSequence := core.Sequence{}
				wholePattern := core.S("1C")
				for repetition := 0; repetition < repetitions*pattern.wholeNotes; repetition++ {
					wholeSequence = wholeSequence.SequenceJoin(wholePattern)
				}
				wholeDuration := core.PlayAt(wholeDevice, core.NoCondition, wholeSequence, bpm, begin).Sub(begin)
				ideal := time.Duration(math.Round(float64(repetitions*pattern.wholeNotes) * 240 * float64(time.Second) / bpm))
				t.Logf("duration=%s, ideal=%s, tempo error=%s, subdivision error=%s", actual, ideal, actual-ideal, actual-wholeDuration)
				if actual != wholeDuration || actual != ideal {
					t.Errorf("duration=%s, equivalent whole notes=%s, ideal=%s", actual, wholeDuration, ideal)
				}
				if got := large.DurationAt(bpm); got != actual {
					t.Errorf("reported duration=%s, playback duration=%s", got, actual)
				}
				var latest time.Time
				timeline.EventsDo(func(event core.TimelineEvent, when time.Time) {
					if when.Before(begin) || when.After(end) {
						t.Errorf("event outside playback interval: %s", when.Sub(begin))
					}
					latest = when
				})
				if pattern.name != "rests" && !latest.Equal(end) {
					t.Errorf("last event at %s, playback ends at %s", latest.Sub(begin), actual)
				}
			})
		}
	}
}

func TestPlaybackClockAcrossSequencesAndTempoChanges(t *testing.T) {
	const repetitions = 1024
	begin := time.Now().Add(time.Hour)
	clock := core.NewPlaybackClock(begin, 123)
	device := NewOutputDevice(1, nil, 1, core.NewTimeline())
	pattern := core.S("32C 16.D")
	large := core.Sequence{}
	for repetition := 0; repetition < repetitions; repetition++ {
		device.PlayWithClock(core.NoCondition, pattern, &clock)
		large = large.SequenceJoin(pattern)
	}
	want := begin.Add(time.Duration(math.Round(repetitions * 0.125 * 240 * float64(time.Second) / 123)))
	wholeDevice := NewOutputDevice(1, nil, 1, core.NewTimeline())
	if got := core.PlayAt(wholeDevice, core.NoCondition, large, 123, begin); !got.Equal(want) || !clock.Time().Equal(got) {
		t.Errorf("partitioned playback=%s, single sequence=%s, want=%s", clock.Time().Sub(begin), got.Sub(begin), want.Sub(begin))
	}
	clock.SetBPM(137)
	if !clock.Time().Equal(want) {
		t.Fatal("changing tempo moved the current boundary")
	}
	for repetition := 0; repetition < repetitions; repetition++ {
		device.PlayWithClock(core.NoCondition, pattern, &clock)
	}
	want = want.Add(time.Duration(math.Round(repetitions * 0.125 * 240 * float64(time.Second) / 137)))
	if !clock.Time().Equal(want) {
		t.Errorf("tempo-changed playback=%s, want=%s", clock.Time().Sub(begin), want.Sub(begin))
	}
}

func TestPlaybackMixedDurations(t *testing.T) {
	fixed, err := core.NewMIDI(core.On(73*time.Millisecond), core.On(62), core.On(80)).ToNote()
	if err != nil {
		t.Fatal(err)
	}
	short, err := core.NewMIDI(core.On(41*time.Millisecond), core.On(65), core.On(80)).ToNote()
	if err != nil {
		t.Fatal(err)
	}
	sequence := core.Sequence{Notes: [][]core.Note{
		{core.N("32C")},
		{fixed},
		{fixed.ToRest()},
		{core.N("16E").WithTiedNote(fixed)},
		{fixed.WithTiedNote(core.N("32C"))},
		{fixed, short},
	}}
	begin := time.Now().Add(time.Hour)
	const bpm = 127
	offset := func(fraction float64, milliseconds int) time.Duration {
		return time.Duration(math.Round(fraction*240*float64(time.Second)/bpm)) + time.Duration(milliseconds)*time.Millisecond
	}
	timeline := core.NewTimeline()
	device := NewOutputDevice(1, nil, 1, timeline)
	end := core.PlayAt(device, core.NoCondition, sequence, bpm, begin)
	if got, want := end.Sub(begin), offset(0.125, 333); got != want {
		t.Errorf("playback duration=%s, want=%s", got, want)
	}
	if got := sequence.DurationAt(bpm); got != end.Sub(begin) {
		t.Errorf("reported duration=%s, playback duration=%s", got, end.Sub(begin))
	}
	expected := []struct {
		number     int
		start, end time.Duration
	}{
		{60, 0, offset(1.0/32, 0)},
		{62, offset(1.0/32, 0), offset(1.0/32, 73)},
		{64, offset(1.0/32, 146), offset(3.0/32, 219)},
		{62, offset(3.0/32, 219), offset(4.0/32, 292)},
		{65, offset(4.0/32, 292), offset(4.0/32, 333)},
		{62, offset(4.0/32, 292), offset(4.0/32, 365)},
	}
	events := timeline.NoteEvents()
	if len(events) != len(expected) {
		t.Fatalf("got %d note events, want %d", len(events), len(expected))
	}
	for index, event := range events {
		want := expected[index]
		if event.Number != want.number || event.Start.Sub(begin) != want.start || event.End.Sub(begin) != want.end {
			t.Errorf("event %d: number=%d start=%s end=%s, want number=%d start=%s end=%s", index, event.Number, event.Start.Sub(begin), event.End.Sub(begin), want.number, want.start, want.end)
		}
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

func Test_canCombineEvent(t *testing.T) {
	type args struct {
		notes []core.Note
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"one", args{[]core.Note{core.MustParseNote("c")}}, true},
		{"cd", args{[]core.Note{core.MustParseNote("c"), core.MustParseNote("d")}}, true},
		{"cd+", args{[]core.Note{core.MustParseNote("c"), core.MustParseNote("d+")}}, false},
		{"cd+", args{[]core.Note{core.MustParseNote(".c"), core.MustParseNote("d")}}, false},
		{"cd+", args{[]core.Note{core.MustParseNote(".c#-"), core.MustParseNote(".d-")}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canCombineEvent(tt.args.notes); got != tt.want {
				t.Errorf("canCombineEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
