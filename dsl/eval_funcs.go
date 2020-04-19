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
	Title         string
	Description   string
	Prefix        string // for autocomplete
	Sample        string
	ControlsAudio bool
	Func          interface{}
}

func EvalFunctions(storage VariableStorage) map[string]Function {
	eval := map[string]Function{}
	eval["chord"] = Function{
		Title:       "Chord creator",
		Description: "create a Chord",
		Prefix:      "cho",
		Sample:      `chord('${1:note}')`,
		Func: func(chord string) interface{} {
			c, err := melrose.ParseChord(chord)
			if err != nil {
				notify.Print(notify.Errorf("%v", err))
				return nil
			}
			return c
		}}

	eval["pitch"] = Function{
		Title:       "Pitch modifier",
		Description: "change the pitch with a delta of semitones",
		Prefix:      "pit",
		Sample:      `pitch(${1:semitones},${2:sequenceable})`,
		Func: func(semitones, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot pitch (%T) %v", m, m))
				return nil
			}
			return melrose.Pitch{Target: s, Semitones: getValueable(semitones)}
		}}

	eval["reverse"] = Function{
		Title:       "Reverse modifier",
		Description: "reverse the (groups of) notes in a sequence",
		Prefix:      "rev",
		Sample:      `reverse(${1:sequenceable})`,
		Func: func(m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot reverse (%T) %v", m, m))
				return nil
			}
			return melrose.Reverse{Target: s}
		}}

	eval["repeat"] = Function{
		Title:       "Repeat modifier",
		Description: "repeat the musical object a number of times",
		Prefix:      "rep",
		Sample:      `repeat(${1:times},${2:sequenceable})`,
		Func: func(howMany int, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot repeat (%T) %v", m, m))
				return nil
			}
			return melrose.Repeat{Target: s, Times: howMany}
		}}

	eval["join"] = Function{
		Title:       "Join modifier",
		Description: "join two or more musical objects",
		Prefix:      "joi",
		Sample:      `join(${1:first},${2:second})`,
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
		Title:         "Beats Per Minute",
		Description:   "set the Beats Per Minute [1..300], default is 120",
		ControlsAudio: true,
		Prefix:        "bpm",
		Sample:        `bpm(${1:beats-per-minute})`,
		Func: func(f float64) interface{} {
			melrose.CurrentDevice().SetBeatsPerMinute(f)
			return nil
		}}

	eval["velocity"] = Function{
		Title:         "controls softness",
		Description:   "set the base velocity [1..127], default is 70",
		ControlsAudio: true,
		Prefix:        "vel",
		Sample:        `velocity(${1:velocity})`,
		Func: func(i int) interface{} {
			// TODO check range
			melrose.CurrentDevice().SetBaseVelocity(i)
			return nil
		}}

	eval["echo"] = Function{
		Title:         "the notes being played",
		Description:   "Echo the notes being played (default is true)",
		ControlsAudio: true,
		Prefix:        "ech",
		Sample:        `echo(${0:true|false})`,
		Func: func(on bool) interface{} {
			melrose.CurrentDevice().SetEchoNotes(on)
			return nil
		}}

	eval["sequence"] = Function{
		Title:       "Sequence creator",
		Description: "create a Sequence from a string of notes",
		Prefix:      "seq",
		Sample:      `sequence('${1:space-separated-notes}')`,
		Func: func(s string) interface{} {
			sq, err := melrose.ParseSequence(s)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return sq
		}}

	eval["note"] = Function{
		Title:       "Note creator",
		Prefix:      "no",
		Description: "Note, e.g. C 2G#5. =",
		Sample:      `note('${1:letter}')`,
		Func: func(s string) interface{} {
			n, err := melrose.ParseNote(s)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return n
		}}

	eval["play"] = Function{
		Title:         "Player (foreground)",
		Description:   "play musical objects such as Note,Chord,Sequence,...",
		ControlsAudio: true,
		Prefix:        "pla",
		Sample:        `play(${1:sequenceable})`,
		Func: func(playables ...interface{}) interface{} {
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					melrose.CurrentDevice().Play(s)
				} else {
					notify.Print(notify.Warningf("cannot play (%T) %v", p, p))
				}
			}
			return nil
		}}

	eval["go"] = Function{
		Title:         "Player (background)",
		Description:   "play all musical objects in parallel",
		ControlsAudio: true,
		Prefix:        "go",
		Sample:        `go(${1:sequenceable})`,
		Func: func(playables ...interface{}) interface{} {
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					go melrose.CurrentDevice().Play(s)
				}
			}
			return nil
		}}

	eval["serial"] = Function{
		Title:       "Serial modifier",
		Description: "serialise any parallelisation of notes in a musical object",
		Prefix:      "ser",
		Sample:      `serial(${1:sequenceable})`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot serial (%T) %v", value, value))
				return nil
			} else {
				return melrose.Serial{Target: s}
			}
		}}

	eval["record"] = Function{
		Title:         "Recorder",
		Description:   "creates a recorded sequence of notes from device ID and stop after T seconds of inactivity",
		ControlsAudio: true,
		Prefix:        "rec",
		Sample:        `record(${1:input-device-id},${1:seconds-inactivity})`,
		Func: func(deviceID int, secondsInactivity int) interface{} {
			seq, err := melrose.CurrentDevice().Record(deviceID, time.Duration(secondsInactivity)*time.Second)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return seq
		}}
	eval["undynamic"] = Function{
		Title:       "Undo Dynamic modifier",
		Description: "undynamic all the notes in a musical object",
		Prefix:      "und",
		Sample:      `undynamic(${1:sequenceable})`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot undynamic (%T) %v", value, value))
				return nil
			} else {
				return melrose.Undynamic{Target: s}
			}
		}}

	eval["flatten"] = Function{
		Title:       "Flatten modifier",
		Description: "flatten all operations on a musical object to a new sequence",
		Prefix:      "flat",
		Sample:      `flatten(${1:sequenceable})`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot flatten (%T) %v", value, value))
				return nil
			} else {
				return s.S()
			}
		}}

	eval["parallel"] = Function{
		Title:       "Parallel modifier",
		Description: "create a new sequence in which all notes of a musical object will be played in parallel",
		Prefix:      "par",
		Sample:      `parallel(${1:sequenceable})`,
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
		Title:       "Loop creator",
		Description: "create a new loop",
		Prefix:      "loo",
		Sample:      `loop(${1:sequenceable}) // stop(${2:variablename})`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot loop (%T) %v", value, value))
				return nil
			} else {
				return &melrose.Loop{Target: s}
			}
		}}
	eval["run"] = Function{
		Title:         "Loop runner",
		Description:   "start loop(s). Ignore if it was running.",
		ControlsAudio: true,
		Prefix:        "run",
		Sample:        `run(${1:loop})`,
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
		Title:         "Loop stopper",
		Description:   "stop running loop(s). Ignore if it was stopped.",
		ControlsAudio: true,
		Prefix:        "sto",
		Sample:        `stop(${1:loop-or-empty})`,
		Func: func(vars ...variable) interface{} {
			if len(vars) == 0 {
				StopAllLoops(storage)
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
		Title:         "MIDI channel modifier",
		Description:   "select a MIDI channel, must be in [0..16]",
		ControlsAudio: true,
		Prefix:        "chan",
		Sample:        `channel(${1:number},${2:sequenceable})`,
		Func: func(midiChannel, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot decorate with channel (%T) %v", m, m))
				return nil
			}
			return melrose.ChannelSelector{Target: s, Number: getValueable(midiChannel)}
		}}
	eval["interval"] = Function{
		Title:         "Integer Interval creator",
		Description:   "create an integer repeating interval (from,to,by)",
		ControlsAudio: true,
		Prefix:        "int",
		Sample:        `interval(${1:from},${2:to},${3:by})`,
		Func: func(from, to, by interface{}) *melrose.Interval {
			return melrose.NewInterval(melrose.ToValueable(from), melrose.ToValueable(to), melrose.ToValueable(by))
		}}
	eval["indexmap"] = Function{
		Title:         "Integer Index Map modifier",
		Description:   "create a Mapper of Notes by index (1-based)",
		ControlsAudio: true,
		Prefix:        "ind",
		Sample:        `indexmap('${0:space-separated-1-based-indices}',${1:sequenceable})`,
		Func: func(indices string, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot create index mapper on (%T) %v", m, m))
				return nil
			}
			return melrose.NewIndexMapper(s, indices)
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
