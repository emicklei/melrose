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
func (s SyncPlay) _Stop(ctx core.Context) error {
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
	for _, each := range s.playables {
		val := each.Value()
		if ply, ok := val.(core.Playable); ok {
			_ = ply.Play(ctx, time.Now())
		} else {
			if seq, ok := val.(core.Sequenceable); ok {
				_ = ctx.Device().Play(core.NoCondition, seq, ctx.Control().BPM(), time.Now())
			}
		}
	}
	return nil
}
