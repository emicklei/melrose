package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/emicklei/melrose/notify"
)

type NoteEvent struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Number   int       `json:"number"`
	Velocity int       `json:"velocity"`
}

func (e NoteEvent) WithEnd(end time.Time) NoteEvent {
	return NoteEvent{Start: e.Start, Number: e.Number, Velocity: e.Velocity, End: end}
}

func NotesEventsToFile(events []NoteEvent, name string) {
	out, err := os.Create(name)
	if err != nil {
		notify.Errorf("NotesEventsToFile:%v", err)
		return
	}
	defer out.Close()
	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	if err := enc.Encode(events); err != nil {
		notify.Errorf("NotesEventsToFile:%v", err)
	}
}

func NoteEventsFromFile(name string) (list []NoteEvent) {
	in, err := os.Open(name)
	if err != nil {
		notify.Errorf("NoteEventsFromFile:%v", err)
		return
	}
	defer in.Close()
	dec := json.NewDecoder(in)
	if err := dec.Decode(&list); err != nil {
		notify.Errorf("NoteEventsFromFile:%v", err)
	}
	return
}

func NoteEventsToPeriods(events []NoteEvent) (list []NotePeriod) {
	for _, each := range events {
		list = append(list, NotePeriod{
			startMs:  each.Start.UnixMilli(),
			endMs:    each.End.UnixMilli(),
			number:   each.Number,
			velocity: each.Velocity,
		})
	}
	return
}

type NoteEventStatistics struct {
	Start, End      time.Time
	Lowest, Highest int
}

func NoteStatistics(list []NoteEvent) (stats NoteEventStatistics) {
	if len(list) == 0 {
		return stats
	}
	stats.Start = list[0].Start
	stats.End = list[len(list)-1].End
	stats.Lowest = 127
	for _, each := range list {
		if each.Number < stats.Lowest {
			stats.Lowest = each.Number
		}
		if each.Number > stats.Highest {
			stats.Highest = each.Number
		}
	}
	return
}

const eventTimeFormat = "03:04:05.000000000"

func (s NoteEventStatistics) String() string {
	return fmt.Sprintf("start=%s,end=%s,low=%d,high=%d", s.Start.Format(eventTimeFormat), s.End.Format(eventTimeFormat), s.Lowest, s.Highest)
}

type NoteChange struct {
	isOn     bool
	note     int64
	velocity int64
}

func (n NoteChange) Number() int {
	return int(n.note)
}

func (n NoteChange) Velocity() int {
	return int(n.velocity)
}

func (n NoteChange) IsOn() bool {
	return n.isOn
}

func (n NoteChange) NoteChangesDo(block func(NoteChange)) { block(n) }

func NewNoteChange(isOn bool, midiNr int64, velocity int64) NoteChange {
	return NoteChange{
		isOn:     isOn,
		note:     midiNr,
		velocity: velocity,
	}
}

func (n NoteChange) Handle(tim *Timeline, when time.Time) {
	// NOP
}
