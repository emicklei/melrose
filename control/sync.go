package control

import (
	"bytes"
	"fmt"
	"time"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/op"
)

type SyncPlay struct {
	playables []core.Valueable
}

func NewSyncPlay(list []core.Valueable) SyncPlay {
	return SyncPlay{playables: list}
}

func (s SyncPlay) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "sync(")
	for i, each := range s.playables {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, "%s", core.Storex(each))
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Play implements Playable
func (s SyncPlay) Play(ctx core.Context, at time.Time) error {
	return s.Evaluate(ctx)
}

// Stop implements Playable
func (s SyncPlay) Stop(ctx core.Context) error {
	for _, each := range s.playables {
		val := each.Value()
		if ply, ok := val.(core.Stoppable); ok {
			_ = ply.Stop(ctx)
		}
	}
	return nil
}

func (s SyncPlay) S() core.Sequence {
	l := []core.Sequenceable{}
	for _, each := range s.playables {
		val := each.Value()
		if s, ok := val.(core.Sequenceable); ok {
			l = append(l, s)
		}
	}
	return (op.Merge{Target: l}).S()
}

// TODO
// IsPlaying implements Playable
// func (s SyncPlay) IsPlaying() bool {
// 	for _, each := range s.playables {
// 		val := each.Value()
// 		if ply, ok := val.(core.Stoppable); ok {
// 			if ply.IsPlaying() {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func (s SyncPlay) Evaluate(ctx core.Context) error {
	// if the first is a Loop at start playing the others on the NextPlayAt
	// if the first is not a Loop then start playing the others now
	if len(s.playables) == 0 {
		return nil
	}
	first := s.playables[0].Value()
	playfirst, ok := first.(core.Playable)
	if !ok {
		// play all right now as sequenceables
		for _, each := range s.playables {
			val := each.Value()
			if ply, ok := val.(core.Sequenceable); ok {
				_ = ctx.Device().Play(core.NoCondition, ply, ctx.Control().BPM(), time.Now())
			}
		}
		return nil
	}
	loopfirst, ok := playfirst.(*core.Loop)
	if !ok {
		// play all right now
		for _, each := range s.playables {
			val := each.Value()
			if ply, ok := val.(core.Playable); ok {
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
		if ply, ok := val.(core.Playable); ok {
			_ = ply.Play(ctx, next)
		}
	}
	return nil
}
