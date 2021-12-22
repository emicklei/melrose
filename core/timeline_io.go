package core

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"time"
)

type StorableNoteChange struct {
	When     int64 `json:"when"`
	IsOn     bool  `json:"ison"`
	Note     int64 `json:"note"`
	Velocity int64 `json:"velocity"`
}

func (t *Timeline) toStorableNoteChanges() (changes []StorableNoteChange) {
	t.EventsDo(func(each TimelineEvent, when time.Time) {
		change, ok := each.(NoteChange)
		if !ok {
			return
		}
		store := StorableNoteChange{
			When:     when.UnixNano(),
			IsOn:     change.isOn,
			Note:     change.note,
			Velocity: change.velocity,
		}
		changes = append(changes, store)
	})
	return
}

func (t *Timeline) ToFile(name string) {
	out, _ := os.Create(name)
	defer out.Close()
	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	if err := enc.Encode(t.toStorableNoteChanges()); err != nil {
		log.Println(err)
	}
}

type NotePeriod struct {
	startMs, endMs int64
	number         int
	velocity       int
}

func (p NotePeriod) Start() time.Time {
	return time.Unix(0, p.startMs*1e6)
}

func (p NotePeriod) End() time.Time {
	return time.Unix(0, p.endMs*1e6)
}

func (p NotePeriod) Number() int { return p.number }

func (p NotePeriod) Velocity() int { return p.velocity }

func (p NotePeriod) Note(bpm float64) Note {
	// TODO assume duration is <= whole note
	sixteenth := int64(math.Round(4 * 60 * 1000 / bpm / 16))
	times := (p.endMs - p.startMs) / sixteenth
	fraction := float32(times) * 0.0625
	n, _ := MIDItoNote(fraction, p.number, p.velocity)
	return n
}

func (p NotePeriod) Quantized(bpm float64) NotePeriod {
	// snap the start to a multiple of 16th note duration for bpm
	// snap the length too
	sixteenth := int64(math.Round(4 * 60 * 1000 / bpm / 16))
	startMs := nearest(p.startMs, sixteenth)
	endMs := nearest(p.endMs, sixteenth)
	return NotePeriod{startMs: startMs, endMs: endMs, number: p.number, velocity: p.velocity}
}

func nearest(value int64, delta int64) int64 {
	return delta * int64(math.RoundToEven((float64(value) / float64(delta))))
}

func ConvertToNotePeriods(changes []StorableNoteChange) (events []NotePeriod) {
	noteOn := map[int64]StorableNoteChange{} // which note started when
	var begin int64 = 0
	for _, each := range changes {
		if each.IsOn {
			noteOn[each.Note] = each
		} else {
			on, ok := noteOn[each.Note]
			if !ok {
				continue
			}
			delete(noteOn, each.Note)
			if begin == 0 {
				begin = on.When
			}
			start := (on.When - begin) / 1e6
			event := NotePeriod{
				startMs:  start,
				endMs:    (each.When - begin) / 1e6,
				number:   int(each.Note),
				velocity: int(on.Velocity),
			}
			events = append(events, event)
		}
	}
	return
}
