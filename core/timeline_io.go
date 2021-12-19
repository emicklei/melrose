package core

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type StorableNoteChange struct {
	When     int64 `json:"when"`
	IsOn     bool  `json:"onoff"`
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

func ConvertToNoteEvents(changes []StorableNoteChange) (events []NoteEvent) {
	noteOn := map[int64]StorableNoteChange{} // which note started when
	for _, each := range changes {
		if each.IsOn {
			noteOn[each.Note] = each
		} else {
			on, ok := noteOn[each.Note]
			if !ok {
				continue
			}
			delete(noteOn, each.Note)
			event := NoteEvent{
				Start:    time.Unix(0, on.When),
				End:      time.Unix(0, each.When),
				Number:   int(each.Note),
				Velocity: int(on.Velocity),
			}
			events = append(events, event)
		}
	}
	return
}
