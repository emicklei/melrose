package core

import (
	"math"
	"time"

	"github.com/emicklei/melrose/notify"
)

// Beatmaster is a LoopController
type Beatmaster struct {
	context         Context
	beating         bool
	bpmChanges      chan float64
	ticker          *time.Ticker
	done            chan bool
	schedule        *BeatSchedule
	beats           int64   // monotonic increasing number, starting at 0
	biab            int64   // current number of beats in a bar
	bpm             float64 // current beats per minute
	settingNotifier func(LoopController)
}

func NewBeatmaster(ctx Context, bpm float64) *Beatmaster {
	return &Beatmaster{
		context:    ctx,
		beating:    false,
		done:       make(chan bool),
		bpmChanges: make(chan float64),
		schedule:   NewBeatSchedule(),
		beats:      0,
		biab:       4,
		bpm:        bpm}
}

func (b *Beatmaster) Reset() {
	b.Stop()
	// drain
	go func() {
		select {
		case <-b.bpmChanges:
		default:
		}
	}()
	b.schedule.Reset()
	b.Start()
}

func (b *Beatmaster) BPM() float64 {
	return b.bpm
}

func (b *Beatmaster) BIAB() int {
	return int(b.biab)
}

func (b *Beatmaster) BeatsAndBars() (int64, int64) {
	return b.beats, b.beats / b.biab
}

func (b *Beatmaster) SettingNotifier(handler func(LoopController)) {
	b.settingNotifier = handler
}

// StartLoop will schedule the start of a Loop at the next bar, unless the master is not started.
func (b *Beatmaster) StartLoop(l *Loop) {
	if !b.beating {
		return
	}
	if l == nil || l.IsRunning() {
		return
	}
	b.schedule.Schedule(b.beatsAtNextBar(), func(when time.Time) {
		l.Start(b.context.Device())
	})
}

// Plan is part of LoopControl
// bars is zero-based
func (b *Beatmaster) Plan(bars int64, seq Sequenceable) {
	atBeats := b.beatsAtNextBar() + (b.biab * bars)
	if IsDebug() {
		notify.Debugf("beat.schedule at: %d put: %s bars: %.2f", atBeats, Storex(seq), seq.S().Bars(int(b.biab)))
	}
	b.schedule.Schedule(atBeats, func(when time.Time) {
		b.context.Device().Play(seq, b.bpm, when)
	})
}

func (b *Beatmaster) beatsAtNextBar() int64 {
	if b.beats%b.biab == 0 {
		return b.beats
	}
	return (b.beats/b.biab + 1) * b.biab
}

// EndLoop will schedule the stop of a Loop at the next bar, unless the master is not started.
func (b *Beatmaster) EndLoop(l *Loop) {
	if !b.beating {
		return
	}
	if l == nil || !l.IsRunning() {
		return
	}
	b.schedule.Schedule(b.beatsAtNextBar(), func(when time.Time) {
		l.Stop()
	})
}

// SetBPM will change the beats per minute at the next bar, unless the master is not started.
func (b *Beatmaster) SetBPM(bpm float64) {
	if !b.beating {
		b.bpm = bpm
		b.notifySettingChanged()
		return
	}
	if b.bpm == bpm {
		return
	}
	go func() { b.bpmChanges <- bpm }()
}

// TODO move checks to SetBIAB in control
// SetBIAB will change the beats per bar, unless the master is not started.
func (b *Beatmaster) SetBIAB(biab int) {
	if !b.beating {
		b.biab = int64(biab)
		return
	}
	if b.biab == int64(biab) {
		return
	}
	b.biab = int64(biab)
	b.notifySettingChanged()
}

func (b *Beatmaster) notifySettingChanged() {
	if b.settingNotifier == nil {
		return
	}
	b.settingNotifier(b)
}

func (b *Beatmaster) Start() {
	if b.beating {
		return
	}
	b.notifySettingChanged()
	b.beats = 0
	b.ticker = time.NewTicker(beatTickerDuration(b.bpm))
	b.beating = true
	go func() {
		if IsDebug() {
			notify.Debugf("core.beatmaster: status=started bpm=%v tick=%v", b.bpm, beatTickerDuration(b.bpm))
		}
		for {
			if b.beats%b.biab == 0 {
				// on a bar
				// abort ?
				select {
				case <-b.done:
					return
				// only change BPM on a bar
				case bpm := <-b.bpmChanges:
					if IsDebug() {
						notify.Debugf("core.beatmaster: bpm=%v tick=%v", bpm, beatTickerDuration(bpm))
					}
					b.bpm = bpm
					b.notifySettingChanged()
					b.ticker.Stop()
					b.ticker = time.NewTicker(beatTickerDuration(bpm))
				default:
				}
			}
			// in between bars
			select {
			case <-b.done:
				return
			case now := <-b.ticker.C:
				if b.schedule.IsEmpty() {
					b.beats = 0
				} else {
					actions := b.schedule.Unschedule(b.beats)
					for _, each := range actions {
						each(now)
					}
					b.beats++
				}
			}
		}
	}()
}

func beatTickerDuration(bpm float64) time.Duration {
	return time.Duration(int(math.Round(float64(60*1000)/bpm))) * time.Millisecond
}

// Stop will stop the beats. Any Loops will continue to run.
func (b *Beatmaster) Stop() {
	if !b.beating {
		return
	}
	b.beating = false
	b.ticker.Stop()
	b.done <- true
	if IsDebug() {
		notify.Debugf("core.beatmaster: stopped")
	}
}

// NoLooper is a Beatmaster that does not loop
var NoLooper = zeroBeat{}

type zeroBeat struct{}

func (s zeroBeat) StartLoop(l *Loop)                            {}
func (s zeroBeat) EndLoop(l *Loop)                              {}
func (s zeroBeat) Start()                                       {}
func (s zeroBeat) Stop()                                        {}
func (s zeroBeat) Reset()                                       {}
func (s zeroBeat) SetBPM(bpm float64)                           {}
func (s zeroBeat) BPM() float64                                 { return 120.0 }
func (s zeroBeat) SetBIAB(biab int)                             {}
func (s zeroBeat) BIAB() int                                    { return 4 }
func (s zeroBeat) BeatsAndBars() (int64, int64)                 { return 0, 0 }
func (s zeroBeat) Plan(bars int64, seq Sequenceable)            {}
func (s zeroBeat) SettingNotifier(handler func(LoopController)) {}
