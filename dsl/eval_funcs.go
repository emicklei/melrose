package dsl

import (
	"strings"
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
)

// Syntax tells what language version this package is supporting.
const Syntax = "1.0" // major,minor

func IsCompatibleSyntax(s string) bool {
	if len(s) == 0 {
		// ignore syntax ; you are on your own
		return true
	}
	mm := strings.Split(Syntax, ".")
	ss := strings.Split(s, ".")
	return mm[0] == ss[0] && ss[1] <= mm[1]
}

type Function struct {
	Description   string
	Aliasses      string // space separated keywords
	Sample        string
	ControlsAudio bool
	Func          interface{}
}

func EvalFunctions(varStore *VariableStore) map[string]Function {
	eval := map[string]Function{}
	eval["chord"] = Function{
		Description: "create a Chord",
		Sample:      `chord('')`,
		Func: func(chord string) interface{} {
			c, err := melrose.ParseChord(chord)
			if err != nil {
				notify.Print(notify.Errorf("%v", err))
				return nil
			}
			return c
		}}

	eval["pitch"] = Function{
		Description: "change the pitch with a delta of semitones",
		Sample:      `pitch(,)`,
		Func: func(semitones, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot pitch (%T) %v", m, m))
				return nil
			}
			return melrose.Pitch{Target: s, Semitones: getValueable(semitones)}
		}}

	eval["reverse"] = Function{
		Description: "reverse the (groups of) notes in a sequence",
		Sample:      `reverse()`,
		Func: func(m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot reverse (%T) %v", m, m))
				return nil
			}
			return melrose.Reverse{Target: s}
		}}

	eval["repeat"] = Function{
		Description: "repeat the musical object a number of times",
		Sample:      `repeat(2,)`,
		Func: func(howMany int, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot repeat (%T) %v", m, m))
				return nil
			}
			return melrose.Repeat{Target: s, Times: howMany}
		}}

	eval["join"] = Function{
		Description: "join two or more musical objects",
		Sample:      `join(,)`,
		Func: func(playables ...interface{}) interface{} { // Note: return type cannot be EvaluationResult
			joined := []melrose.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot join (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return melrose.Join{List: joined}
		}}

	eval["bpm"] = Function{
		Description:   "get or set the Beats Per Minute value [1..300], default is 120",
		ControlsAudio: true,
		Sample:        `bpm(180)`,
		Func: func(f ...float64) interface{} {
			if len(f) == 0 {
				return melrose.CurrentDevice().BeatsPerMinute()
			}
			melrose.CurrentDevice().SetBeatsPerMinute(f[0])
			return nil
		}}

	eval["sequence"] = Function{
		Description: "create a Sequence from a string of notes",
		Sample:      `sequence('')`,
		Aliasses:    "seq",
		Func: func(s string) interface{} {
			sq, err := melrose.ParseSequence(s)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return sq
		}}

	eval["note"] = Function{
		Description: "create a Note from a string",
		Sample:      `note('')`,
		Func: func(s string) interface{} {
			n, err := melrose.ParseNote(s)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return n
		}}

	eval["play"] = Function{
		Description:   "play a musical object such as Note,Chord,Sequence,...",
		ControlsAudio: true,
		Sample:        `play()`,
		Func: func(playables ...interface{}) interface{} {
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					melrose.CurrentDevice().Play(s, true)
				} else {
					notify.Print(notify.Warningf("cannot play (%T) %v", p, p))
				}
			}
			return nil
		}}

	eval["go"] = Function{
		Description:   "play all musical objects in parallel",
		ControlsAudio: true,
		Sample:        `go()`,
		Func: func(playables ...interface{}) interface{} {
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					go melrose.CurrentDevice().Play(s, false)
				}
			}
			return nil
		}}

	eval["serial"] = Function{
		Description: "serialise any parallelisation of notes in a musical object",
		Sample:      `serial()`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot serial (%T) %v", value, value))
				return nil
			} else {
				return melrose.Serial{Target: s}
			}
		}}

	eval["record"] = Function{
		Description:   "creates a recorded sequence of notes from device ID and stop after T seconds of inactivity",
		ControlsAudio: true,
		Sample:        `record(,)`,
		Func: func(deviceID int, secondsInactivity int) interface{} {
			seq, err := melrose.CurrentDevice().Record(deviceID, time.Duration(secondsInactivity)*time.Second)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return seq
		}}
	eval["undynamic"] = Function{
		Description: "undynamic all the notes in a musical object",
		Sample:      `undynamic()`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot undynamic (%T) %v", value, value))
				return nil
			} else {
				return melrose.Undynamic{Target: s}
			}
		}}

	eval["flatten"] = Function{
		Description: "flatten all operations on a musical object to a new sequence",
		Sample:      `flatten()`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot flatten (%T) %v", value, value))
				return nil
			} else {
				return s.S()
			}
		}}

	eval["parallel"] = Function{
		Description: "create a new sequence in which all notes of a musical object will be played in parallel",
		Sample:      `parallel()`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot parallel (%T) %v", value, value))
				return nil
			} else {
				return melrose.Parallel{Target: s}
			}
		}}
	// BEGIN Loop and control
	eval["loop"] = Function{
		Description: "create a new loop",
		Sample:      `loop(s)`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot loop (%T) %v", value, value))
				return nil
			} else {
				return &melrose.Loop{Target: s}
			}
		}}
	eval["run"] = Function{
		Description:   "start loop(s). Ignore if it was running.",
		ControlsAudio: true,
		Sample:        `run(l)`,
		Func: func(vars ...variable) interface{} {
			for _, each := range vars {
				l, ok := each.Value().(*melrose.Loop)
				if !ok {
					notify.Print(notify.Warningf("cannot start (%T) %v", l, l))
					continue
				}
				l.Start(melrose.CurrentDevice())
				notify.Print(notify.Infof("started loop: %s", each.Name))
			}
			return nil
		}}
	eval["stop"] = Function{
		Description:   "stop running loop(s). Ignore if it was stopped.",
		ControlsAudio: true,
		Sample:        `stop(l)`,
		Func: func(vars ...variable) interface{} {
			if len(vars) == 0 {
				StopAllLoops(varStore)
				return nil
			}
			for _, each := range vars {
				l, ok := each.Value().(*melrose.Loop)
				if !ok {
					notify.Print(notify.Warningf("cannot stop (%T) %v", l, l))
					continue
				}
				notify.Print(notify.Infof("stopping loop: %s", each.Name))
				l.Stop()
			}
			return nil
		}}
	// END Loop and control
	eval["channel"] = Function{
		Description:   "select a MIDI channel, must be in [1..16]",
		ControlsAudio: true,
		Sample:        `channel()`,
		Func: func(midiChannel, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot decorate with channel (%T) %v", m, m))
				return nil
			}
			return melrose.ChannelSelector{Target: s, Number: getValueable(midiChannel)}
		}}
	eval["interval"] = Function{
		Description:   "create an integer interval [from..to] with a by.",
		ControlsAudio: true,
		Sample:        `interval()`,
		Func: func(from, to, by interface{}) *melrose.Interval {
			return melrose.NewInterval(melrose.ToValueable(from), melrose.ToValueable(to), melrose.ToValueable(by))
		}}
	return eval
}

func getSequenceable(v interface{}) (melrose.Sequenceable, bool) {
	if s, ok := v.(melrose.Sequenceable); ok {
		return s, ok
	}
	return nil, false
}

func getValueable(val interface{}) melrose.Valueable {
	if v, ok := val.(melrose.Valueable); ok {
		return v
	}
	return melrose.On(val)
}
