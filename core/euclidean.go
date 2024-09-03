package core

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type Euclidean struct {
	Steps    HasValue
	Beats    HasValue
	Rotation HasValue
	Playback HasValue
}

// Play is part of Playable
func (e *Euclidean) Play(ctx Context, at time.Time) error {
	steps := getInt(e.Steps, false)
	beats := getInt(e.Beats, false)
	rotation := getInt(e.Rotation, false)
	playback, ok := e.Playback.Value().(Sequenceable)
	if !ok {
		return fmt.Errorf("playback is not a sequenceable")
	}
	toggles := euclideanRhythm(steps, beats, rotation)
	bpm := ctx.Control().BPM()
	moment := at
	dt := WholeNoteDuration(bpm) / time.Duration(steps)
	for _, each := range toggles {
		if each {
			ctx.Device().Play(NoCondition, playback, bpm, moment)
		}
		moment = moment.Add(dt)
	}
	return nil
}

// Handle is part of TimelineEvent
func (e *Euclidean) Handle(tim *Timeline, when time.Time) {}

func (e *Euclidean) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "euclidean(%s,%s,%s,%s)", Storex(e.Steps), Storex(e.Beats), Storex(e.Rotation), Storex(e.Playback))
	return b.String()
}

func (e *Euclidean) Evaluate(ctx Context) error {
	// copied from Loop
	e.Play(ctx, time.Now())
	return nil
}

// Inspect is part of Inspectable
func (e *Euclidean) Inspect(i Inspection) {
	steps := getInt(e.Steps, false)
	beats := getInt(e.Beats, false)
	rotation := getInt(e.Rotation, false)
	toggles := euclideanRhythm(steps, beats, rotation)
	b := new(strings.Builder)
	for _, each := range toggles {
		if each {
			b.WriteString("!")
		} else {
			b.WriteString(".")
		}
	}
	i.Properties["pattern"] = b.String()
	i.Properties["steps"] = e.Steps
	i.Properties["beats"] = e.Beats
	i.Properties["rotation"] = e.Rotation
	i.Properties["playback"] = e.Playback
}

// https://github.com/computermusicdesign/euclidean-rhythm/blob/master/max-example/euclidSimple.js
func euclideanRhythm(steps, pulses, rotation int) []bool {
	rhythm := make([]bool, 0, steps)
	bucket := 0

	for i := 0; i < steps; i++ {
		bucket += pulses
		if bucket >= steps {
			bucket -= steps
			rhythm = append(rhythm, true)
		} else {
			rhythm = append(rhythm, false)
		}
	}

	return rotatedRhythm(rhythm, rotation+1)
}

func rotatedRhythm(input []bool, rotate int) []bool {
	output := make([]bool, len(input))
	val := len(input) - rotate
	for i := 0; i < len(input); i++ {
		j := (i + val) % len(input)
		if j < 0 {
			j *= -1
		}
		output[i] = input[j]
	}
	return output
}
