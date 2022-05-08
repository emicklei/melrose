package core

import (
	"math"
	"time"
)

func WholeNoteDuration(bpm float64) time.Duration {
	return time.Duration(int(math.Round(4*60*1000/bpm))) * time.Millisecond
}

func FractionToDurationParts(f float64) (fraction float32, dotted bool) {
	type duration struct {
		fraction float32
		dotted   bool
		actual   float64
	}
	durations := []duration{
		{1.0, true, 1.5},
		{1.0, false, 1.0},
		{0.5, true, 0.75},
		{0.5, false, 0.5},
		{0.25, true, 0.375},
		{0.25, false, 0.25},
		{0.125, true, 0.1875},
		{0.125, false, 0.125},
		{0.0625, true, 0.09375},
		{0.0625, false, 0.0625},
	}
	hitDistance := 2.0
	hit := durations[0]
	for _, each := range durations {
		if distance := abs64(each.actual - f); distance <= hitDistance {
			hit = each
			hitDistance = distance
		}
	}
	return hit.fraction, hit.dotted
}

func abs64(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
