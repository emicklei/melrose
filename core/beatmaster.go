package core

import (
	"fmt"
	"math"
	"time"

	"github.com/emicklei/melrose/notify"
)

// Beatmaster is a LoopController
type Beatmaster struct {
	context    Context
	beating    bool
	bpmChanges chan float64
	ticker     *time.Ticker
	done       chan bool
	schedule   *BeatSchedule
	beats      int64   // monotonic increasing number, starting at 0
	biab       int64   // current number of beats in a bar
	bpm        float64 // current beats per minute

	// verbose will produce cryptic logging of the behavior
	// . = a quarter tick
	// | = a bar
	// * = a loop was started
	// x = a loop was stopped
	verbose bool
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
		bpm:        bpm,
		verbose:    false}
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

func (b *Beatmaster) Verbose(v bool) {
	b.verbose = v
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

// StartLoop will schedule the start of a Loop at the next bar, unless the master is not started.
func (b *Beatmaster) StartLoop(l *Loop) {
	if !b.beating {
		return
	}
	if l == nil || l.IsRunning() {
		return
	}
	b.schedule.Schedule(b.beatsAndNextBar(), func(beats int64) {
		l.Start(b.context.Device())
	})
}

// Plan is part of LoopControl
func (b *Beatmaster) Plan(bars int64, beats int64, seq Sequenceable) {
	if !b.beating {
		return
	}
	atBeats := b.beatsAndNextBar()
	atBeats += (b.biab * bars)
	atBeats += beats
	b.schedule.Schedule(atBeats, func(beats int64) {
		b.context.Device().Play(seq, b.bpm, time.Now())
	})
}

func (b *Beatmaster) beatsAndNextBar() int64 {
	if b.beats%b.biab == 0 {
		return b.beats
	}
	return (b.beats * b.biab / b.biab) + 1
}

// EndLoop will schedule the stop of a Loop at the next bar, unless the master is not started.
func (b *Beatmaster) EndLoop(l *Loop) {
	if !b.beating {
		return
	}
	if l == nil || !l.IsRunning() {
		return
	}
	b.schedule.Schedule(b.beatsAndNextBar(), func(b int64) {
		l.Stop()
	})
}

// SetBPM will change the beats per minute at the next bar, unless the master is not started.
func (b *Beatmaster) SetBPM(bpm float64) {
	if !b.beating {
		b.bpm = bpm
		return
	}
	if b.bpm == bpm {
		return
	}
	if bpm < 1.0 {
		notify.Print(notify.Warningf("bpm [%.1f] must be in [1..300], setting to [%d]", bpm, 1))
		bpm = 1.0
	}
	if bpm > 300 {
		notify.Print(notify.Warningf("bpm [%.1f] must be in [1..300], setting to [%d]", bpm, 300))
		bpm = 300.0
	}
	b.bpm = bpm
	go func() { b.bpmChanges <- bpm }()
}

// SetBIAB will change the beats per bar, unless the master is not started.
func (b *Beatmaster) SetBIAB(biab int) {
	if !b.beating {
		b.biab = int64(biab)
		return
	}
	if b.biab == int64(biab) {
		return
	}
	if biab < 1 {
		notify.Print(notify.Warningf("biab [%d] must be in [1..6], setting to [%d]", biab, 1))
		biab = 1
	}
	if biab > 6 {
		notify.Print(notify.Warningf("biab [%d] must be in [1..6], setting to [%d]", biab, 6))
		biab = 6
	}
	b.biab = int64(biab)
}

func (b *Beatmaster) Start() {
	if b.beating {
		return
	}
	b.beats = 0
	b.ticker = time.NewTicker(tickerDuration(b.bpm))
	b.beating = true
	if b.verbose {
		fmt.Println("beatmaster started")
	}
	go func() {
		for {
			if b.beats%b.biab == 0 {
				// on a bar
				// abort ?
				select {
				case <-b.done:
					return
				// only change BPM on a bar
				case bpm := <-b.bpmChanges:
					if b.verbose {
						fmt.Println("bpm:", bpm)
					}
					b.bpm = bpm
					b.ticker.Stop()
					b.ticker = time.NewTicker(tickerDuration(bpm))
				default:
				}
			}
			// in between bars
			select {
			case <-b.done:
				return
			case <-b.ticker.C:
				if b.schedule.IsEmpty() {
					b.beats = 0
				} else {
					actions := b.schedule.Unschedule(b.beats)
					for _, each := range actions {
						each(b.beats)
					}
					b.beats++
				}
			}
		}
	}()
}

func tickerDuration(bpm float64) time.Duration {
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
	if b.verbose {
		fmt.Println("\nbeatmaster stopped")
	}
}

// NoLooper is a Beatmaster that does not loop
var NoLooper = zeroBeat{}

type zeroBeat struct{}

func (s zeroBeat) StartLoop(l *Loop)                              {}
func (s zeroBeat) EndLoop(l *Loop)                                {}
func (s zeroBeat) Start()                                         {}
func (s zeroBeat) Stop()                                          {}
func (s zeroBeat) Reset()                                         {}
func (s zeroBeat) SetBPM(bpm float64)                             {}
func (s zeroBeat) BPM() float64                                   { return 120.0 }
func (s zeroBeat) SetBIAB(biab int)                               {}
func (s zeroBeat) BIAB() int                                      { return 4 }
func (s zeroBeat) BeatsAndBars() (int64, int64)                   { return 0, 0 }
func (s zeroBeat) Plan(bars int64, beats int64, seq Sequenceable) {}
