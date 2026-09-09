package main

import (
	"fmt"
	"os"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/ui/img"
	"github.com/fogleman/gg"
	"gitlab.com/gomidi/midi/v2/smf"
)

// midiNoteEvent captures a note on/off with its absolute tick position within a track.
type midiNoteEvent struct {
	absTicks int64
	isOn     bool
	key      uint8
	velocity uint8
}

func main() {
	s, err := smf.ReadFile("fur-elise.mid")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read midi file:", err)
		os.Exit(1)
	}

	var bpmSum float64
	var bpmCount int
	var noteEvents []midiNoteEvent
	for trackNo, track := range s.Tracks {
		var absTicks int64
		for _, ev := range track {
			absTicks += int64(ev.Delta)

			var channel, key, velocity uint8
			var bpm float64
			if ev.Message.GetNoteOn(&channel, &key, &velocity) {
				fmt.Printf("track %d tick %d: note-on  channel=%d key=%d velocity=%d\n", trackNo, absTicks, channel, key, velocity)
				noteEvents = append(noteEvents, midiNoteEvent{absTicks: absTicks, isOn: true, key: key, velocity: velocity})
			} else if ev.Message.GetNoteOff(&channel, &key, &velocity) {
				fmt.Printf("track %d tick %d: note-off channel=%d key=%d velocity=%d\n", trackNo, absTicks, channel, key, velocity)
				noteEvents = append(noteEvents, midiNoteEvent{absTicks: absTicks, isOn: false, key: key, velocity: velocity})
			} else if ev.Message.GetMetaTempo(&bpm) {
				fmt.Printf("track %d tick %d: bpm=%.2f\n", trackNo, absTicks, bpm)
				bpmSum += bpm
				bpmCount++
			}
		}
	}

	if bpmCount == 0 {
		return
	}
	averageBPM := bpmSum / float64(bpmCount)
	fmt.Printf("average bpm=%.2f\n", averageBPM)

	metricTicks, ok := s.TimeFormat.(smf.MetricTicks)
	if !ok {
		fmt.Fprintln(os.Stderr, "cannot schedule timeline: unsupported time format")
		return
	}

	timeline := core.NewTimeline()
	start := time.Now().Add(500 * time.Millisecond) // give the scheduler a moment before the first event
	for _, each := range noteEvents {
		when := start.Add(metricTicks.Duration(averageBPM, uint32(each.absTicks)))
		timeline.Schedule(core.NewNoteChange(each.isOn, int64(each.key), int64(each.velocity)), when)
	}

	gc := gg.NewContext(2000, 150)
	nv := img.NotesView{Events: timeline.NoteEvents(), BPM: averageBPM}
	nv.DrawOn(gc)
	if err := gc.SavePNG("fur-elise.png"); err != nil {
		fmt.Fprintln(os.Stderr, "failed to save image:", err)
		os.Exit(1)
	}

	periods := timeline.BuildNotePeriods()
	builder := core.NewSequenceBuilder(periods, averageBPM)
	seq := builder.Build()
	fmt.Printf("sequence built: %+v\n", seq)
}
