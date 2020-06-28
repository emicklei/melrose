package melrose

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

type Track struct {
	Title   string
	Channel int
	Content map[int]Sequenceable // bar -> musical object
}

func NewTrack(title string, channel int) *Track {
	return &Track{
		Title:   title,
		Channel: channel,
		Content: map[int]Sequenceable{},
	}
}

// S is part of Sequenceable
func (t *Track) S() Sequence {
	// TODO
	if one, ok := t.Content[1]; ok {
		return one.S()
	}
	return EmptySequence
}

func (t *Track) Play(ctx Context) error {
	now := time.Now()
	bpm := ctx.Control().BPM()
	whole := WholeNoteDuration(bpm)
	for bars, each := range t.Content {
		cs := NewChannelSelector(each, On(t.Channel))
		ctx.Device().Play(cs, bpm, now.Add(time.Duration(bars-1)*whole))
	}
	return nil
}

func (t *Track) Inspect(i Inspection) {
	i.Properties["title"] = t.Title
	i.Properties["channel"] = t.Channel
	i.Properties["pieces"] = len(t.Content)
}

// Add adds a SequenceOnTrack or a Sequence at bar 1.
func (t *Track) Add(seq interface{}) {
	if at, ok := seq.(SequenceOnTrack); ok {
		t.Content[Int(at.Bar)] = at.Target
		return
	}
	if s, ok := seq.(Sequenceable); ok {
		t.Content[1] = s
	}
}

// Storex implements Storable
func (t *Track) Storex() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "track('%s',%d", t.Title, t.Channel)
	for k, v := range t.Content {
		fmt.Fprintf(&buf, ",")
		sont := NewSequenceOnTrack(On(k), v)
		fmt.Fprintf(&buf, sont.Storex())
	}
	fmt.Fprintf(&buf, ")")
	return buf.String()
}

type SequenceOnTrack struct {
	Bar    Valueable
	Target Sequenceable
}

func NewSequenceOnTrack(bar Valueable, seq Sequenceable) SequenceOnTrack {
	return SequenceOnTrack{Bar: bar, Target: seq}
}

func (s SequenceOnTrack) S() Sequence {
	return s.Target.S()
}

// Storex implements Storable
func (s SequenceOnTrack) Storex() string {
	if st, ok := s.Target.(Storable); ok {
		return fmt.Sprintf("put(%v,%s)", s.Bar, st.Storex())
	}
	return ""
}

type MultiTrack struct {
	Tracks []Valueable
}

// Storex implements Storable
func (m MultiTrack) Storex() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "multi(")
	for i, each := range m.Tracks {
		if i > 0 {
			fmt.Fprintf(&buf, ",")
		}
		if t, ok := each.(Storable); ok {
			fmt.Fprintf(&buf, "%s", t.Storex())
		} else {
			fmt.Fprintf(&buf, "%v", each)
		}
	}
	fmt.Fprintf(&buf, ")")
	return buf.String()
}

// Play is part of Playable
func (m MultiTrack) Play(ctx Context) error {
	for _, each := range m.Tracks {
		if track, ok := each.Value().(*Track); ok {
			for bar, seq := range track.Content {
				ch := ChannelSelector{Number: On(track.Channel), Target: seq}
				ctx.Control().Plan(int64(bar-1), int64(0), ch)
			}
		} else {
			// TODO
			log.Println("not a track")
		}
	}
	return nil
}
