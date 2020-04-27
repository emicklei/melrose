package melrose

import (
	"fmt"
	"math"
	"time"
)

// Beatmaster is a LoopController
type Beatmaster struct {
	beating        bool
	ticker         *time.Ticker
	done           chan bool
	loopStartQueue chan *Loop
	loopStopQueue  chan *Loop
	bpmChanges     chan float64
	biabChanges    chan int64
	beats          int64   // monotonic increasing number, starting at 0
	biab           int64   // current number of beats in a bar
	bpm            float64 // current beats per minute
	verbose        bool    // if true log beats and bars
}

func NewBeatmaster(bpm float64) *Beatmaster {
	return &Beatmaster{
		beating:        false,
		ticker:         time.NewTicker(tickerDuration(bpm)),
		done:           make(chan bool),
		loopStartQueue: make(chan *Loop),
		loopStopQueue:  make(chan *Loop),
		bpmChanges:     make(chan float64),
		biabChanges:    make(chan int64),
		beats:          0,
		biab:           4,
		bpm:            bpm,
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

func (b *Beatmaster) BPM() float64 {
	return b.bpm
}

func (b *Beatmaster) BIAB() int {
	return int(b.biab)
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

// SetBPM will change the beats per minute at the next bar, unless the master is not started.
func (b *Beatmaster) SetBPM(bpm float64) {
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

// SetBIAB will change the beats per bar, unless the master is not started.
func (b *Beatmaster) SetBIAB(biab int) {
	if !b.beating {
		return
	}
	if biab < 1 || biab > 6 { // TODO what is the max?
		return
	}
	go func() {
		b.biabChanges <- int64(biab)
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
				if b.beats%b.biab == 0 {
					if b.verbose {
						fmt.Print("!")
					}

					select {
					// abort ?
					case <-b.done:
						return
					// only change BIAB on a bar
					case biab := <-b.biabChanges:
						if b.verbose {
							fmt.Println("biab:", biab)
						}
						b.biab = biab
					default:
					}

					select {
					// abort ?
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
							l.Start(Context().AudioDevice)
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

var NoLooper = zeroBeat{}

type zeroBeat struct{}

func (s zeroBeat) Begin(l *Loop)      {}
func (s zeroBeat) End(l *Loop)        {}
func (s zeroBeat) Start()             {}
func (s zeroBeat) Stop()              {}
func (s zeroBeat) SetBPM(bpm float64) {}
func (s zeroBeat) BPM() float64       { return 120.0 }
func (s zeroBeat) SetBIAB(biab int)   {}
func (s zeroBeat) BIAB() int          { return 4 }
