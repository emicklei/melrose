package core

import (
	"bytes"
	"fmt"
	"time"
)

type SyncPlay struct {
	playables []Valueable
}

func NewSyncPlay(list []Valueable) SyncPlay {
	return SyncPlay{playables: list}
}

func (s SyncPlay) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "sync(")
	for i, each := range s.playables {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, "%s", Storex(each))
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Play implements Playable
func (s SyncPlay) Play(ctx Context, at time.Time) error {
	return s.Evaluate(ctx)
}

// Stop implements Playable
func (s SyncPlay) Stop(ctx Context) error {
	for _, each := range s.playables {
		val := each.Value()
		if ply, ok := val.(Playable); ok {
			_ = ply.Stop(ctx)
		}
	}
	return nil
}

func (s SyncPlay) Evaluate(ctx Context) error {
	// if the first is a Loop at start playing the others on the NextPlayAt
	// if the first is not a Loop then start playing the others now
	if len(s.playables) == 0 {
		return nil
	}
	first := s.playables[0].Value()
	playfirst, ok := first.(Playable)
	if !ok {
		// play all right now as sequenceables
		for _, each := range s.playables {
			val := each.Value()
			if ply, ok := val.(Sequenceable); ok {
				_ = ctx.Device().Play(NoCondition, ply, ctx.Control().BPM(), time.Now())
			}
		}
		return nil
	}
	loopfirst, ok := playfirst.(*Loop)
	if !ok {
		// play all right now
		for _, each := range s.playables {
			val := each.Value()
			if ply, ok := val.(Playable); ok {
				_ = ply.Play(ctx, time.Now())
			}
		}
		return nil
	}
	// handle first loop
	next := loopfirst.NextPlayAt()
	// running?
	if next.IsZero() {
		next = time.Now()
	}
	for _, each := range s.playables {
		val := each.Value()
		if ply, ok := val.(Playable); ok {
			_ = ply.Play(ctx, next)
		}
	}
	return nil
}
