package core

import (
	"bytes"
	"fmt"
	"time"

	"github.com/emicklei/melrose/notify"
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

func (t *Track) Play(ctx Context, now time.Time) error {
	bpm := ctx.Control().BPM()
	biab := ctx.Control().BIAB()
	whole := WholeNoteDuration(bpm)
	for bars, each := range t.Content {
		cs := NewChannelSelector(each, On(t.Channel))
		offset := int64((bars-1)*biab) * whole.Nanoseconds() / 4
		when := now.Add(time.Duration(time.Duration(offset)))
		if IsDebug() {
			notify.Debugf("core.track title=%s channel=%d bar=%d, biab=%d, bpm=%.2f time=%s", t.Title, t.Channel, bars, biab, bpm, when.Format("04:05.000"))
		}
		ctx.Device().Play(NoCondition, cs, bpm, when)
	}
	return nil
}

func (t *Track) Inspect(i Inspection) {
	i.Properties["channel"] = t.Channel
	i.Properties["pieces"] = len(t.Content)
}

// Add adds a SequenceOnTrack
func (t *Track) Add(seq SequenceOnTrack) {
	b := Int(seq.Bar)
	t.Content[b] = seq.Target
}

// Storex implements Storable
func (t *Track) Storex() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "track('%s',%d", t.Title, t.Channel)
	for k, v := range t.Content {
		fmt.Fprint(&buf, ",")
		sont := NewSequenceOnTrack(On(k), v) // TODO
		fmt.Fprint(&buf, sont.Storex())
	}
	fmt.Fprintf(&buf, ")")
	return buf.String()
}

type SequenceOnTrack struct {
	Bar    HasValue
	Target Sequenceable
}

func NewSequenceOnTrack(bar HasValue, seq Sequenceable) SequenceOnTrack {
	return SequenceOnTrack{Bar: bar, Target: seq}
}

func (s SequenceOnTrack) S() Sequence {
	return s.Target.S()
}

// Storex implements Storable
func (s SequenceOnTrack) Storex() string {
	if st, ok := s.Target.(Storable); ok {
		return fmt.Sprintf("onbar(%v,%s)", s.Bar, st.Storex())
	}
	return ""
}

type MultiTrack struct {
	Tracks []HasValue
}

// Storex implements Storable
func (m MultiTrack) Storex() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "multitrack(")
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
func (m MultiTrack) Play(ctx Context, at time.Time) error {
	// because all tracks must be synchronized, we first stop the beatmaster
	// then schedule all tracks
	// then start the beatmaster again.
	ctx.Control().Stop()
	for _, each := range m.Tracks {
		if track, ok := each.Value().(*Track); ok {
			for bar, seq := range track.Content {
				ch := NewChannelSelector(seq, On(track.Channel))
				ctx.Control().Plan(int64(bar-1), ch)
			}
		} else {
			// TODO
			notify.NewErrorf("not a track:%v", each)
		}
	}
	ctx.Control().Start()
	return nil
}
