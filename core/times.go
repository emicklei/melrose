package core

import (
	"math"
	"time"
)

func WholeNoteDuration(bpm float64) time.Duration {
	return fractionDuration(1, bpm)
}

func fractionDuration(fraction, bpm float64) time.Duration {
	return time.Duration(math.Round(fraction * 240 * float64(time.Second) / bpm))
}

const musicalTicksPerWhole = 1 << 32

type PlaybackClock struct {
	origin time.Time
	bpm    float64
	ticks  int64
	fixed  time.Duration
}

func NewPlaybackClock(begin time.Time, bpm float64) PlaybackClock {
	return PlaybackClock{origin: begin, bpm: bpm}
}

func (c PlaybackClock) Duration() time.Duration {
	return fractionDuration(float64(c.ticks)/musicalTicksPerWhole, c.bpm) + c.fixed
}

func (c PlaybackClock) Time() time.Time {
	return c.origin.Add(c.Duration())
}

func (c *PlaybackClock) SetBPM(bpm float64) {
	if c.bpm == bpm {
		return
	}
	c.origin = c.Time()
	c.bpm = bpm
	c.ticks = 0
	c.fixed = 0
}

func (c PlaybackClock) After(note Note) PlaybackClock {
	if note.duration > 0 {
		c.fixed += note.duration
	} else {
		fraction := float64(note.fraction)
		if note.Dotted {
			fraction *= 1.5
		}
		c.ticks += int64(math.Round(fraction * musicalTicksPerWhole))
	}
	for _, tied := range note.tied {
		c = c.After(tied)
	}
	return c
}

func (c PlaybackClock) AfterGroup(notes []Note) PlaybackClock {
	if len(notes) == 0 {
		return c
	}
	earliest := c.After(notes[0])
	for _, note := range notes[1:] {
		end := c.After(note)
		if end.Duration() < earliest.Duration() {
			earliest = end
		}
	}
	return earliest
}

func (c PlaybackClock) AfterSequence(sequence Sequence) PlaybackClock {
	for _, group := range sequence.Notes {
		c = c.AfterGroup(group)
	}
	return c
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
