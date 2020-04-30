package dsl

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/op"
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

func EvalFunctions(storage VariableStorage, control melrose.LoopController) map[string]Function {
	eval := map[string]Function{}
	eval["chord"] = Function{
		Title:       "Chord creator",
		Description: "create a Chord",
		Prefix:      "cho",
		Alias:       "C",
		Template:    `chord('${1:note}')`,
		Samples: `chord('C#5/m/1')
chord('G/M/2')`,
		IsCore: true,
		Func: func(chord string) interface{} {
			c, err := melrose.ParseChord(chord)
			if err != nil {
				notify.Print(notify.Errorf("%v", err))
				return nil
			}
			return c
		}}

	eval["widechord"] = Function{
		Title:       "Chord widener",
		Description: "make a Chord wider",
		Prefix:      "wid",
		Template:    `widechord('${1:chord}')`,
		IsCore:      false,
		Func: func(chord interface{}) op.WideChord {
			// chord is a melrose.Chord
			// or a variable with a chord
			return op.WideChord{Target: getValueable(chord)}
		}}

	eval["pitch"] = Function{
		Title:       "Pitch modifier",
		Description: "change the pitch with a delta of semitones",
		Prefix:      "pit",
		Alias:       "Pi",
		Template:    `pitch(${1:semitones},${2:sequenceable})`,
		Samples: `pitch(-1,sequence('C D E'))
p = interval(-4,4,1)
pitch(p,note('C'))`,
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
		Func: func(playables ...interface{}) interface{} {
			joined := []melrose.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot join (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.Join{Target: joined}
		}}

	eval["bpm"] = Function{
		Title:         "Beats Per Minute",
		Description:   "set the Beats Per Minute [1..300]; default is 120",
		ControlsAudio: true,
		Prefix:        "bpm",
		Template:      `bpm(${1:beats-per-minute})`,
		Func: func(f float64) interface{} {
			melrose.Context().LoopControl.SetBPM(f)
			return nil
		}}

	eval["biab"] = Function{
		Title:         "Beats in a Bar",
		Description:   "set the Beats in a Bar [1..6]; default is 4",
		ControlsAudio: true,
		Prefix:        "biab",
		Template:      `biab(${1:beats-in-a-bar})`,
		Func: func(i int) interface{} {
			if i < 1 || i > 6 {
				notify.Print(notify.Warningf("invalid beats-in-a-bar [1..6], %d = ", i))
				return nil
			}
			melrose.Context().LoopControl.SetBIAB(i)
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
			melrose.Context().AudioDevice.SetBaseVelocity(i)
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
			melrose.Context().AudioDevice.SetEchoNotes(on)
			return nil
		}}

	eval["sequence"] = Function{
		Title:       "Sequence creator",
		Description: "create a Sequence from (space separated) notes",
		Prefix:      "seq",
		Alias:       "S",
		Template:    `sequence('${1:space-separated-notes}')`,
		Samples: `sequence('C D E')
sequence('[C D E]')`,
		IsCore: true,
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
		Description: "create a Scale using a starting Note and type indicator (Major,minor)",
		Prefix:      "sc",
		Template:    `scale('${1:letter}')`,
		IsCore:      true,
		Samples:     `scale('E/m') // => E F G A B C5 D5`,
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
			moment := time.Now()
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					moment = melrose.Context().AudioDevice.Play(s, melrose.Context().LoopControl.BPM(), moment)
				} else {
					notify.Print(notify.Warningf("cannot play (%T) %v", p, p))
				}
			}
			// wait until the play is completed before allowing a new one to.
			// add a bit of time to allow the previous play to finish all its notes.
			// time.Sleep(moment.Sub(time.Now()) + (50 * time.Millisecond))
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
			moment := time.Now()
			for _, p := range playables {
				if s, ok := getSequenceable(p); ok {
					melrose.Context().AudioDevice.Play(s, melrose.Context().LoopControl.BPM(), moment)
				} else {
					notify.Print(notify.Warningf("cannot go (%T) %v", p, p))
				}
			}
			return nil
		}}

	eval["serial"] = Function{
		Title:       "Serial modifier",
		Description: "serialise any parallelisation of notes in one or more musical objects",
		Prefix:      "ser",
		Template:    `serial(${1:sequenceable})`,
		IsComposer:  true,
		Samples: `serial(chord('E')) // => E G B
serial(sequence('[C D]'),note('E')) // => C D E`,
		Func: func(playables ...interface{}) interface{} {
			joined := []melrose.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot serial (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return melrose.Serial{Target: joined}
		}}

	eval["octave"] = Function{
		Title:       "Octave modifier",
		Description: "changes the pitch of notes by steps of 12 semitones",
		Prefix:      "oct",
		Template:    `octave(${1:offet},${2:sequenceable})`,
		IsComposer:  true,
		Samples:     `octave(1,sequence('C D')) // => C5 D5`,
		Func: func(scalarOrVar interface{}, playables ...interface{}) interface{} {
			joined := []melrose.Sequenceable{}
			for _, p := range playables {
				if s, ok := getSequenceable(p); !ok {
					notify.Print(notify.Warningf("cannot octave (%T) %v", p, p))
					return nil
				} else {
					joined = append(joined, s)
				}
			}
			return op.Octave{Target: joined, Offset: melrose.ToValueable(scalarOrVar)}
		}}

	eval["record"] = Function{
		Title:         "Recorder",
		Description:   "creates a recorded sequence of notes from a MIDI device",
		ControlsAudio: true,
		Prefix:        "rec",
		Template:      `record(${1:input-device-id},${1:seconds-inactivity})`,
		Samples: `r = record(1,5) // record notes played on device ID=1 and stop recording after 5 seconds
s = r.Sequence()`,
		Func: func(deviceID int, secondsInactivity int) interface{} {
			seq, err := melrose.Context().AudioDevice.Record(deviceID, time.Duration(secondsInactivity)*time.Second)
			if err != nil {
				notify.Print(notify.Error(err))
				return nil
			}
			return seq
		}}
	eval["undynamic"] = Function{
		Title:       "Undo dynamic modifier",
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
		Samples:     `flatten(sequence('[C E G] B')) // => C E G B`,
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
		Description: "create a new sequence in which all notes of a musical object are synched in time",
		Prefix:      "par",
		Alias:       "Pa",
		Template:    `parallel(${1:sequenceable})`,
		Samples:     `parallel(sequence('C D E')) // => [C D E]`,
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
		Title:         "Loop creator ; must be assigned to a variable",
		Description:   "create a new loop",
		ControlsAudio: true,
		Prefix:        "loo",
		Alias:         "L",
		Template:      `lp_${1:object} = loop(${1:object}) // end(lp_${1:object})`,
		Samples: `cb = sequence('C D E F G A B')
lp_cb = loop(cb)`,
		Func: func(value interface{}) interface{} {
			if s, ok := getSequenceable(value); !ok {
				notify.Print(notify.Warningf("cannot loop (%T) %v", value, value))
				return nil
			} else {
				return melrose.NewLoop(s)
			}
		}}
	eval["begin"] = Function{
		Title:         "Loop runner",
		Description:   "begin loop(s). Ignore if it was running.",
		ControlsAudio: true,
		Prefix:        "beg",
		Template:      `begin(${1:loop})`,
		Samples: `l1 = loop(sequence('C D E F G A B'))
end(l1)
begin(l1)`,
		Func: func(vars ...variable) interface{} {
			for _, each := range vars {
				l, ok := each.Value().(*melrose.Loop)
				if !ok {
					notify.Print(notify.Warningf("cannot begin (%T) %v", l, l))
					continue
				}
				control.Begin(l)
				notify.Print(notify.Infof("started loop: %s", each.Name))
			}
			return nil
		}}
	eval["end"] = Function{
		Title:         "Loop terminator",
		Description:   "end running loop(s). Ignore if it was stopped.",
		ControlsAudio: true,
		Template:      `end(${1:loop-or-empty})`,
		Samples: `l1 = loop(sequence('C E G))
end(l1)`,
		Func: func(vars ...variable) interface{} {
			if len(vars) == 0 {
				StopAllLoops(storage)
				return nil
			}
			for _, each := range vars {
				l, ok := each.Value().(*melrose.Loop)
				if !ok {
					notify.Print(notify.Warningf("cannot end (%T) %v", l, l))
					continue
				}
				notify.Print(notify.Infof("stopping loop: %s", each.Name))
				control.End(l)
			}
			return nil
		}}
	// END Loop and control
	eval["channel"] = Function{
		Title:         "MIDI channel modifier ; must be a top-level modifier",
		Description:   "select a MIDI channel, must be in [0..16]",
		ControlsAudio: true,
		Prefix:        "chan",
		Alias:         "Ch",
		Template:      `channel(${1:number},${2:sequenceable})`,
		Samples:       `channel(2,sequence('C2 E3') // plays on instrument connected to MIDI channel 2'`,
		Func: func(midiChannel, m interface{}) interface{} {
			s, ok := getSequenceable(m)
			if !ok {
				notify.Print(notify.Warningf("cannot decorate with channel (%T) %v", m, m))
				return nil
			}
			return melrose.ChannelSelector{Target: s, Number: getValueable(midiChannel)}
		}}
	eval["interval"] = Function{
		Title:         "Integer Interval creator. Repeats on default",
		Description:   "create an integer repeating interval (from,to,by)",
		ControlsAudio: false,
		Prefix:        "int",
		Alias:         "I",
		Template:      `interval(${1:from},${2:to},${3:by})`,
		Samples: `i1 = interval(-2,4,1)
l1 = loop(pitch(i1,sequence('C D E F')))`,
		IsComposer: true,
		Func: func(from, to, by interface{}) *melrose.Interval {
			return melrose.NewInterval(melrose.ToValueable(from), melrose.ToValueable(to), melrose.ToValueable(by), melrose.RepeatFromTo)
		}}
	eval["indexmap"] = Function{
		Title:         "Integer Index Map modifier",
		Description:   "create a Mapper of Notes by index (1-based)",
		ControlsAudio: false,
		Prefix:        "ind",
		Alias:         "Im",
		Template:      `indexmap('${1:space-separated-1-based-indices}',${2:sequenceable})`,
		Samples: `s1 = sequence('C D E F G A B')
i1 = indexmap('6 5 4 3 2 1',s1) // => B A G F E D`,
		IsComposer: true,
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
