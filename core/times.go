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

// musicalTicksPerWhole is the tick resolution of a whole note; a power of two so
// that all note fractions (including 1/32 and dotted values) are exact in ticks.
const musicalTicksPerWhole = 1 << 32

// PlaybackClock tracks a position in time relative to an origin.
// Musical positions are accumulated as integer ticks to avoid rounding drift,
// while notes with an absolute duration are accumulated separately.
// It is a value type: After* methods return an advanced copy and leave the receiver unchanged.
type PlaybackClock struct {
	// origin is the wall-clock time at which the current tempo segment started.
	origin time.Time
	// bpm is the tempo used to convert ticks to time.
	bpm float64
	// ticks is the musical position (in 1/musicalTicksPerWhole whole notes) since origin.
	ticks int64
	// fixed is the position contributed by notes with a fixed duration, independent of tempo.
	fixed time.Duration
}

// NewPlaybackClock returns a clock positioned at begin that runs at the given tempo.
func NewPlaybackClock(begin time.Time, bpm float64) PlaybackClock {
	return PlaybackClock{origin: begin, bpm: bpm}
}

// Duration returns the time elapsed since the origin.
func (c PlaybackClock) Duration() time.Duration {
	// Convert the tick total to a whole-note fraction and round once, only on the absolute offset.
	return fractionDuration(float64(c.ticks)/musicalTicksPerWhole, c.bpm) + c.fixed
}

// Time returns the absolute wall-clock time of the current position.
func (c PlaybackClock) Time() time.Time {
	return c.origin.Add(c.Duration())
}

// SetBPM changes the tempo from the current position onwards.
// Time already elapsed is preserved by moving the origin to the current position.
func (c *PlaybackClock) SetBPM(bpm float64) {
	if c.bpm == bpm {
		return
	}
	// Rebase: elapsed time was computed with the old tempo and must not be rescaled.
	c.origin = c.Time()
	c.bpm = bpm
	c.ticks = 0
	c.fixed = 0
}

// After returns a copy of the clock advanced by the length of the note, including any tied notes.
func (c PlaybackClock) After(note Note) PlaybackClock {
	if note.duration > 0 {
		// A fixed duration is not scaled by tempo.
		c.fixed += note.duration
	} else {
		fraction := float64(note.fraction)
		if note.Dotted {
			// A dot extends the note by half of its length.
			fraction *= 1.5
		}
		c.ticks += int64(math.Round(fraction * musicalTicksPerWhole))
	}
	// Tied notes sound as one, so their lengths are added.
	for _, tied := range note.tied {
		c = c.After(tied)
	}
	return c
}

// AfterGroup returns a copy of the clock advanced past a group of simultaneous notes.
// The clock advances by the shortest note in the group, matching MIDI playback
// where the next group starts when the earliest note ends.
func (c PlaybackClock) AfterGroup(notes []Note) PlaybackClock {
	if len(notes) == 0 {
		return c
	}
	earliest := c.After(notes[0])
	for _, note := range notes[1:] {
		// Each note is measured from the same start position c, not from the previous note.
		end := c.After(note)
		if end.Duration() < earliest.Duration() {
			earliest = end
		}
	}
	return earliest
}

// AfterSequence returns a copy of the clock advanced past all note groups of the sequence.
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
		{0.03125, true, 0.046875},
		{0.03125, false, 0.03125},
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
