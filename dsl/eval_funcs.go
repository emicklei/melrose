package dsl

import (
	"fmt"
	"log"
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
	Alias         string // short notation
	Template      string // for autocomplete in VSC
	Samples       string // for doc generation
	ControlsAudio bool
	Tags          string // space separated
	IsCore        bool   // creates a core musical object
	IsComposer    bool   // can decorate a musical object or other decorations
	Func          interface{}
}

func EvalFunctions(storage VariableStorage) map[string]Function {
	eval := map[string]Function{}
	eval["chord"] = Function{
		Title:       "Chord creator",
		Description: "create a Chord",
		Prefix:      "cho",
		Alias:       "C",
		Template:    `chord('${1:note}')`,
		Samples: `chord('C#5/m/1')
chord('G/M/2)`,
		IsCore: true,
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
		Alias:       "Pi",
		Template:    `pitch(${1:semitones},${2:sequenceable})`,
		Samples: `pitch(-1,sequence('C D E'))
pitch(12,note('C'))`,
		IsComposer: true,
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
		Alias:       "Rv",
		Template:    `reverse(${1:sequenceable})`,
		Samples:     `reverse(chord('A'))`,
		IsComposer:  true,
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
		Alias:       "Rp",
		Template:    `repeat(${1:times},${2:sequenceable})`,
		Samples:     `repeat(4,sequence('C D E'))`,
		IsComposer:  true,
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
		Alias:       "J",
		Template:    `join(${1:first},${2:second})`,
		IsComposer:  true,
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
		Description:   "set the Beats Per Minute [1..300]; default is 120",
		ControlsAudio: true,
		Prefix:        "bpm",
		Template:      `bpm(${1:beats-per-minute})`,
		Func: func(f float64) interface{} {
			melrose.CurrentDevice().SetBeatsPerMinute(f)
			return nil
		}}

	eval["velocity"] = Function{
		Title:         "controls softness",
		Description:   "set the base velocity [1..127]; default is 70",
		ControlsAudio: true,
		Prefix:        "vel",
		Template:      `velocity(${1:velocity})`,
		Samples:       `velocity(90)`,
		Func: func(i int) interface{} {
			// TODO check range
			melrose.CurrentDevice().SetBaseVelocity(i)
			return nil
		}}

	eval["echo"] = Function{
		Title:         "the notes being played",
		Description:   "echo the notes being played; default is true",
		ControlsAudio: true,
		Prefix:        "ech",
		Template:      `echo(${0:true|false})`,
		Samples:       `echo(false)`,
		Func: func(on bool) interface{} {
			melrose.CurrentDevice().SetEchoNotes(on)
			return nil
		}}

	eval["sequence"] = Function{
		Title:       "Sequence creator",
		Description: "create a Sequence from (space separated) notes",
		Prefix:      "seq",
		Alias:       "S",
		Template:    `sequence('${1:space-separated-notes}')`,
		Samples:     `sequence('C D E')`,
		IsCore:      true,
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
		Description: "create a Note from the note notation",
		Prefix:      "no",
		Alias:       "N",
		Template:    `note('${1:letter}')`,
		Samples: `note('E')
note('2E#.--')`,
		IsCore: true,
		Func: func(s string) interface{} {
			n, err := melrose.ParseNote(s)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return n
		}}

	eval["scale"] = Function{
		Title:       "Scale creator",
		Prefix:      "sc",
		Description: "",
		Template:    `scale('${1:letter}')`,
		IsCore:      true,
		Samples:     `scale('C#/m')`,
		Func: func(s string) interface{} {
			sc, err := melrose.ParseScale(s)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return sc
		}}

	eval["play"] = Function{
		Title:         "Player (foreground)",
		Description:   "play musical objects such as Note,Chord,Sequence,...",
		ControlsAudio: true,
		Prefix:        "pla",
		Template:      `play(${1:sequenceable})`,
		Samples:       `play(s1,s2,s3) // play s3 after s2 after s1`,
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
		Template:      `go(${1:sequenceable})`,
		Samples:       `go(s1,s1,s3) // play s1 and s2 and s3 simultaneously`,
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
		Template:    `serial(${1:sequenceable})`,
		IsComposer:  true,
		Samples:     `serial(chord('E')) => E G B`,
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
		Template:      `record(${1:input-device-id},${1:seconds-inactivity})`,
		Samples:       `record(1,5) // record notes played on device ID=1 and stop recording after 5 seconds`,
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
		Template:    `undynamic(${1:sequenceable})`,
		IsComposer:  true,
		Samples:     `undynamic('A+ B++ C-- D-') // =>  A B C D`,
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
		Alias:       "F",
		Template:    `flatten(${1:sequenceable})`,
		IsComposer:  true,
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
		Alias:       "Pa",
		Template:    `parallel(${1:sequenceable})`,
		IsComposer:  true,
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
		Title:         "Loop creator",
		Description:   "create a new loop",
		ControlsAudio: true,
		Prefix:        "loo",
		Alias:         "L",
		Template:      `loop(${1:sequenceable}) // stop(${2:variablename})`,
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
		Template:      `run(${1:loop})`,
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
		Template:      `stop(${1:loop-or-empty})`,
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
				// Stop waits for the loop to end so run it in a go-routine
				go l.Stop()
			}
			return nil
		}}
	// END Loop and control
	eval["channel"] = Function{
		Title:         "MIDI channel modifier",
		Description:   "select a MIDI channel, must be in [0..16]",
		ControlsAudio: true,
		Prefix:        "chan",
		Alias:         "Ch",
		Template:      `channel(${1:number},${2:sequenceable})`,
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
		ControlsAudio: false,
		Prefix:        "int",
		Alias:         "I",
		Template:      `interval(${1:from},${2:to},${3:by})`,
		IsComposer:    true,
		Func: func(from, to, by interface{}) *melrose.Interval {
			return melrose.NewInterval(melrose.ToValueable(from), melrose.ToValueable(to), melrose.ToValueable(by))
		}}
	eval["indexmap"] = Function{
		Title:         "Integer Index Map modifier",
		Description:   "create a Mapper of Notes by index (1-based)",
		ControlsAudio: false,
		Prefix:        "ind",
		Alias:         "Im",
		Template:      `indexmap('${1:space-separated-1-based-indices}',${2:sequenceable})`,
		IsComposer:    true,
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

func registerFunction(m map[string]Function, k string, f Function) {
	if dup, ok := m[k]; ok {
		log.Fatal("duplicate function key detected:", dup)
	}
	if len(f.Alias) > 0 {
		if dup, ok := m[f.Alias]; ok {
			log.Fatal("duplicate function alias key detected:", dup)
		}
	}
	m[k] = f
	if len(f.Alias) > 0 {
		// modify title
		f.Title = fmt.Sprintf("%s [%s]", f.Title, f.Alias)
		m[f.Alias] = f
	}
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
