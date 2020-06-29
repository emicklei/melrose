package core

import (
	"math"
	"time"
)

func WholeNoteDuration(bpm float64) time.Duration {
	return time.Duration(int(math.Round(4*60*1000/bpm))) * time.Millisecond
}
