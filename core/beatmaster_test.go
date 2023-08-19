package core

import (
	"testing"
	"time"
)

func Test_tickerDuration(t *testing.T) {
	d := beatTickerDuration(60.0)
	if got, want := d, 1*time.Second; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	d2 := beatTickerDuration(300.0)
	if got, want := d2, 200*time.Millisecond; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	d3 := beatTickerDuration(120.0)
	if got, want := d3, 500*time.Millisecond; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBeatmaster_beatsAtNextBar(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.beats = 0
	b.biab = 3
	if got, want := b.beatsAtNextBar(), int64(0); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	b.beats = 5
	b.biab = 4
	if got, want := b.beatsAtNextBar(), int64(8); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTrackBarTiming(t *testing.T) {
	ctx := PlayContext{}
	b := NewBeatmaster(ctx, 120.0)
	b.SetBPM(100.0)
	b.SetBIAB(3)
	ctx.LoopControl = b

	tr1 := NewTrack("1", 1)
	tr1.Add(NewSequenceOnTrack(On(1), MustParseSequence("c d e")))
	tr2 := NewTrack("2", 1)
	tr1.Add(NewSequenceOnTrack(On(2), MustParseSequence("c")))
	m := MultiTrack{Tracks: []HasValue{On(tr1), On(tr2)}}
	m.Play(ctx, time.Now())
	t.Log(b.schedule.entries)
	_, ok := b.schedule.entries[0]
	if !ok {

		t.Fail()
	}
	_, ok = b.schedule.entries[3]
	if !ok {
		t.Fail()
	}
}
