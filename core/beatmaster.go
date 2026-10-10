package core

import (
	"sync"
	"time"

	"github.com/emicklei/melrose/notify"
)

// Beatmaster is a LoopController
type Beatmaster struct {
	mu              sync.RWMutex // guards beating, beats, biab, bpm and settingNotifier
	context         Context
	beating         bool
	bpmChanges      chan float64
	timer           *time.Timer
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
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.bpm
}

func (b *Beatmaster) BIAB() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return int(b.biab)
}

func (b *Beatmaster) BeatsAndBars() (int64, int64) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.beats, b.beats / b.biab
}

func (b *Beatmaster) SettingNotifier(handler func(LoopController)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.settingNotifier = handler
}

// Plan is part of LoopControl
// bars is zero-based
func (b *Beatmaster) Plan(bars int64, seq Sequenceable) {
	b.mu.RLock()
	biab := b.biab
	atBeats := b.beatsAtNextBar() + (biab * bars)
	b.mu.RUnlock()
	notify.Debugf("beat.schedule at beats: %d put: %s bars: %.2f", atBeats, Storex(seq), seq.S().Bars(int(biab)))

	b.schedule.Schedule(atBeats, func(when time.Time) {
		d := b.context.Device()
		if d != nil { // TODO happens on testing; NEEDSFIX
			PlayAt(d, NoCondition, seq, b.BPM(), when)
		}
	})
}

// beatsAtNextBar requires b.mu to be held.
func (b *Beatmaster) beatsAtNextBar() int64 {
	if b.beats%b.biab == 0 {
		return b.beats
	}
	return (b.beats/b.biab + 1) * b.biab
}

// SetBPM will change the beats per minute at the next bar, unless the master is not started.
func (b *Beatmaster) SetBPM(bpm float64) {
	b.mu.Lock()
	if !b.beating {
		b.bpm = bpm
		b.mu.Unlock()
		b.notifySettingChanged()
		return
	}
	if b.bpm == bpm {
		b.mu.Unlock()
		return
	}
	if b.schedule.IsEmpty() {
		b.bpm = bpm
		b.mu.Unlock()
		b.notifySettingChanged()
		return
	}
	b.mu.Unlock()
	go func() { b.bpmChanges <- bpm }()
}

// TODO move checks to SetBIAB in control
// SetBIAB will change the beats per bar, unless the master is not started.
func (b *Beatmaster) SetBIAB(biab int) {
	b.mu.Lock()
	if !b.beating {
		b.biab = int64(biab)
		b.mu.Unlock()
		return
	}
	if b.biab == int64(biab) {
		b.mu.Unlock()
		return
	}
	b.biab = int64(biab)
	b.mu.Unlock()
	b.notifySettingChanged()
}

// notifySettingChanged calls the handler without holding the lock because it may call back.
func (b *Beatmaster) notifySettingChanged() {
	b.mu.RLock()
	handler := b.settingNotifier
	b.mu.RUnlock()
	if handler == nil {
		return
	}
	handler(b)
}

func (b *Beatmaster) Start() {
	b.mu.RLock()
	alreadyBeating := b.beating
	b.mu.RUnlock()
	if alreadyBeating {
		return
	}
	b.notifySettingChanged()
	b.mu.Lock()
	if b.beating {
		b.mu.Unlock()
		return
	}
	b.beats = 0
	clock := NewPlaybackClock(time.Now(), b.bpm)
	nextBeat := clock.After(Rest4)
	timer := time.NewTimer(time.Until(nextBeat.Time()))
	b.timer = timer
	b.beating = true
	b.mu.Unlock()
	go func() {
		defer timer.Stop()
		if notify.IsDebug() {
			bpm := b.BPM()
			notify.Debugf("core.beatmaster: started bpm=%v tick=%v", bpm, beatTickerDuration(bpm))
		}
		for {
			b.mu.RLock()
			onBar := b.beats%b.biab == 0
			b.mu.RUnlock()
			if onBar {
				// on a bar
				// abort ?
				select {
				case <-b.done:
					return
				// only change BPM on a bar
				case bpm := <-b.bpmChanges:
					if notify.IsDebug() {
						notify.Debugf("core.beatmaster: changed bpm=%v tick=%v", bpm, beatTickerDuration(bpm))
					}
					b.mu.Lock()
					b.bpm = bpm
					b.mu.Unlock()
					b.notifySettingChanged()
				default:
				}
			}
			if bpm := b.BPM(); clock.bpm != bpm {
				clock.SetBPM(bpm)
				nextBeat = clock.After(Rest4)
				timer.Reset(time.Until(nextBeat.Time()))
			}
			// in between bars
			select {
			case <-b.done:
				return
			case <-timer.C:
				clock = nextBeat
				if b.schedule.IsEmpty() {
					b.mu.Lock()
					b.beats = 0
					b.mu.Unlock()
				} else {
					beats, _ := b.BeatsAndBars()
					actions := b.schedule.Unschedule(beats)
					for _, each := range actions {
						each(clock.Time())
					}
					b.mu.Lock()
					b.beats++
					b.mu.Unlock()
				}
				nextBeat = clock.After(Rest4)
				timer.Reset(time.Until(nextBeat.Time()))
			}
		}
	}()
}

func beatTickerDuration(bpm float64) time.Duration {
	return fractionDuration(0.25, bpm)
}

// Stop will stop the beats. Any Loops will continue to run.
func (b *Beatmaster) Stop() {
	b.mu.Lock()
	if !b.beating {
		b.mu.Unlock()
		return
	}
	b.beating = false
	b.mu.Unlock()
	b.done <- true
	if notify.IsDebug() {
		notify.Debugf("core.beatmaster: stopped")
	}
}

// NoLooper is a Beatmaster that does not loop
var NoLooper = zeroBeat{}

type zeroBeat struct{}

func (s zeroBeat) Start()                                       {}
func (s zeroBeat) Stop()                                        {}
func (s zeroBeat) Reset()                                       {}
func (s zeroBeat) SetBPM(bpm float64)                           {}
func (s zeroBeat) BPM() float64                                 { return 120.0 } // same as the default tempo
func (s zeroBeat) SetBIAB(biab int)                             {}
func (s zeroBeat) BIAB() int                                    { return 4 }
func (s zeroBeat) BeatsAndBars() (int64, int64)                 { return 0, 0 }
func (s zeroBeat) Plan(bars int64, seq Sequenceable)            {}
func (s zeroBeat) SettingNotifier(handler func(LoopController)) {}
