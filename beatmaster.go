package melrose

import (
	"fmt"
	"math"
	"time"
)

type Beatmaster struct {
	beating        bool
	ticker         *time.Ticker
	done           chan bool
	loopStartQueue chan *Loop
	loopStopQueue  chan *Loop
	bpmChanges     chan float64
	beats          int64 // monotonic increasing number, starting at 0
	bar            int64 // bar = number of beats on which a new bar starts
	verbose        bool  // if true log beats and bars
}

func NewBeatmaster(bpm float64) *Beatmaster {
	td := tickerDuration(4, bpm)
	return &Beatmaster{
		beating:        false,
		ticker:         time.NewTicker(td),
		done:           make(chan bool),
		loopStartQueue: make(chan *Loop),
		loopStopQueue:  make(chan *Loop),
		bpmChanges:     make(chan float64),
		beats:          0,
		bar:            4,
		verbose:        false}
}

// Verbose will produce cryptic logging of the behavior
// . = a quarter tick
// | = a bar
// * = a loop was started
// x = a loop was stopped
func (b *Beatmaster) Verbose(v bool) {
	b.verbose = v
}

// Begin will schedule the start of a Loop at the next bar, unless the master is not started.
func (b *Beatmaster) Begin(l *Loop) {
	if !b.beating {
		return
	}
	if l == nil || l.IsRunning() {
		return
	}
	go func() {
		b.loopStartQueue <- l
	}()
}

// End will schedule the stop of a Loop at the next bar, unless the master is not started.
func (b *Beatmaster) End(l *Loop) {
	if !b.beating {
		return
	}
	if l == nil || !l.IsRunning() {
		return
	}
	go func() {
		b.loopStopQueue <- l
	}()
}

// BPM will change the beats per minute at the next bar, unless the master is not started.
func (b *Beatmaster) BPM(bpm float64) {
	if !b.beating {
		return
	}
	if bpm < 1.0 || bpm > 300.0 { // TODO what is the max?
		return
	}
	go func() {
		b.bpmChanges <- bpm
	}()
}

func (b *Beatmaster) Start() {
	if b.beating {
		return
	}
	b.beating = true
	if b.verbose {
		fmt.Println("beatmaster started")
	}
	go func() {
		for {
			select {
			// abort ?
			case <-b.done:
				return
			case <-b.ticker.C:
				if b.beats%b.bar == 0 {
					if b.verbose {
						fmt.Print("!")
					}
					select {
					// abort ?
					case <-b.done:
						return
					// do we need to change BPM ?
					case bpm := <-b.bpmChanges:
						b.ticker.Stop()
						b.ticker = time.NewTicker(tickerDuration(b.bar, bpm))
					default:
					}
					// start all pending loops in start Q
					// stop all pending loops in stop Q
					for {
						select {
						// abort ?
						case <-b.done:
							return
						case l := <-b.loopStartQueue:
							if b.verbose {
								fmt.Print("*")
							}
							l.Start(CurrentDevice())
						case l := <-b.loopStopQueue:
							if b.verbose {
								fmt.Print("x")
							}
							l.Stop()
						default:
							goto emptyQueues
						}
					}
				emptyQueues:
				} else {
					if b.verbose {
						fmt.Print(".")
					}
				}
				b.beats++
			}
		}
	}()
}

func tickerDuration(bar int64, bpm float64) time.Duration {
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

var NoLooper = silentLooper{}

type silentLooper struct{}

func (s silentLooper) Begin(l *Loop)   {}
func (s silentLooper) End(l *Loop)     {}
func (s silentLooper) Start()          {}
func (s silentLooper) Stop()           {}
func (s silentLooper) BPM(bpm float64) {}

type UnscheduledLooper struct{}

func (u UnscheduledLooper) Begin(l *Loop) {
	l.Start(CurrentDevice())
}

func (u UnscheduledLooper) End(l *Loop) {
	l.Stop()
}

func (u UnscheduledLooper) Start()      {}
func (u UnscheduledLooper) Stop()       {}
func (u UnscheduledLooper) BPM(float64) {}
