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
		Description: "create a triad Chord with a Note",
		Sample:      `chord('')`,
		Func: func(chord string) FunctionResult {
			c, err := melrose.ParseChord(chord)
			if err != nil {
				return result(nil, notify.Errorf("%v", err))
			}
			return result(c, nil)
		}}

	eval["pitch"] = Function{
		Description: "change the pitch with a delta of semitones",
		Sample:      `pitch(1,)`,
		Func: func(semitones int, m interface{}) FunctionResult {
			s, ok := getSequenceable(m)
			if !ok {
				return result(nil, notify.Warningf("cannot pitch (%T) %v", m, m))
			}
			return result(melrose.Pitch{Target: s, Semitones: semitones}, nil)
		}}

	eval["reverse"] = Function{
		Description: "reverse the (groups of) notes in a sequence",
		Sample:      `reverse()`,
		Func: func(m interface{}) FunctionResult {
			s, ok := getSequenceable(m)
			if !ok {
				return result(nil, notify.Warningf("cannot reverse (%T) %v", m, m))
			}
			return result(melrose.Reverse{Target: s}, nil)
		}}

	eval["repeat"] = Function{
		Description: "repeat the musical object a number of times",
		Sample:      `repeat(2,)`,
		Func: func(howMany int, m interface{}) FunctionResult {
			s, ok := getSequenceable(m)
			if !ok {
				return result(nil, notify.Warningf("cannot repeat (%T) %v", m, m))
			}
			return result(melrose.Repeat{Target: s, Times: howMany}, nil)
		}}

	eval["join"] = Function{
		Description: "join two or more musical objects",
		Sample:      `join(,)`,
		Func: func(playables ...interface{}) interface{} { // Note: return type cannot be EvaluationResult
			joined := []melrose.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					joined = append(joined, s)
				} else {
					return result(nil, notify.Warningf("cannot join (%T) %v", p, p))
				}
			}
			return result(melrose.Join{List: joined}, nil)
		}}

	eval["bpm"] = Function{
		Description:   "get or set the Beats Per Minute value [1..300], default is 120",
		ControlsAudio: true,
		Sample:        `bpm(180)`,
		Func: func(f ...float64) FunctionResult {
			if len(f) == 0 {
				return result(melrose.CurrentDevice().BeatsPerMinute(), nil)
			}
			melrose.CurrentDevice().SetBeatsPerMinute(f[0])
			return result(f[0], nil)
		}}

	eval["sequence"] = Function{
		Description: "create a Sequence from a string of notes",
		Sample:      `sequence('')`,
		Aliasses:    "seq",
		Func: func(s string) FunctionResult {
			n, err := melrose.ParseSequence(s)
			if err != nil {
				return result(nil, notify.Error(err))
			}
			return result(n, nil)
		}}

	eval["note"] = Function{
		Description: "create a Note from a string",
		Sample:      `note('')`,
		Func: func(s string) FunctionResult {
			n, err := melrose.ParseNote(s)
			if err != nil {
				return result(nil, notify.Error(err))
			}
			return result(n, nil)
		}}

	eval["play"] = Function{
		Description:   "play a musical object such as Note,Chord,Sequence,...",
		ControlsAudio: true,
		Sample:        `play()`,
		Func: func(playables ...interface{}) interface{} { // Note: return type cannot be EvaluationResult
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					melrose.CurrentDevice().Play(s, true)
				} else {
					return result(nil, notify.Warningf("cannot play (%T) %v", p, p))
				}
			}
			return result(nil, nil)
		}}

	eval["go"] = Function{
		Description:   "play all musical objects in parallel",
		ControlsAudio: true,
		Sample:        `go()`,
		Func: func(playables ...interface{}) interface{} { // Note: return type cannot be EvaluationResult
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					go melrose.CurrentDevice().Play(s, false)
				}
			}
			return result(nil, nil)
		}}

	eval["serial"] = Function{
		Description: "serialise any parallelisation of notes in a musical object",
		Sample:      `serial()`,
		Func: func(value interface{}) FunctionResult {
			if s, ok := getSequenceable(value); ok {
				return result(melrose.Serial{Target: s}, nil)
			} else {
				return result(nil, notify.Warningf("cannot serial (%T) %v", value, value))
			}
		}}

	eval["record"] = Function{
		Description:   "creates a recorded sequence of notes from device ID and stop after T seconds of inactivity",
		ControlsAudio: true,
		Sample:        `record(,)`,
		Func: func(deviceID int, secondsInactivity int) FunctionResult {
			seq, err := melrose.CurrentDevice().Record(deviceID, time.Duration(secondsInactivity)*time.Second)
			return result(seq, notify.Error(err))
		}}

	eval["undynamic"] = Function{
		Description: "undynamic all the notes in a musical object",
		Sample:      `undynamic()`,
		Func: func(value interface{}) FunctionResult {
			if s, ok := getSequenceable(value); ok {
				return result(melrose.Undynamic{Target: s}, nil)
			} else {
				return result(nil, notify.Warningf("cannot undynamic (%T) %v", value, value))
			}
		}}

	eval["flatten"] = Function{
		Description: "flatten all operations on a musical object (mo) to a new sequence",
		Sample:      `flatten()`,
		Func: func(value interface{}) FunctionResult {
			if s, ok := getSequenceable(value); ok {
				return result(s.S(), nil)
			} else {
				return result(nil, notify.Warningf("cannot flatten (%T) %v", value, value))
			}
		}}

	eval["parallel"] = Function{
		Description: "create a new sequence in which all notes of a musical object will be played in parallel",
		Sample:      `parallel()`,
		Func: func(value interface{}) FunctionResult {
			if s, ok := getSequenceable(value); ok {
				return result(melrose.Parallel{Target: s}, nil)
			} else {
				return result(nil, notify.Warningf("cannot parallel (%T) %v", value, value))
			}
		}}

	return eval
}

func getSequenceable(v interface{}) (melrose.Sequenceable, bool) {
	if s, ok := v.(melrose.Sequenceable); ok {
		return s, ok
	}
	if f, ok := v.(FunctionResult); ok {
		if f.Notification != nil {
			notify.Print(f.Notification)
		}
		return getSequenceable(f.Result)
	}
	return nil, false
}
